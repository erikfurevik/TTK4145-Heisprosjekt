package ElevatorController

import (
	nc "../NetworkController"
	"../config"
	"../elevio"
	"../fsm"
)

/*This is the brain of our local elevator. It has the following resposibilities
	- Distribute and redistribute orders around to elevators
	- Set order lights
	- Handle fault tolerance such as engine failure, loss of network communication
	- Sync the local elevator data with every other external elevator
	- Pass orders to fsm
*/
func MainLogicFunction(Local_ID int, HardwareToControl <-chan elevio.ButtonEvent, UpdateLight chan<- [config.NumElevator]config.Elev,
	LocalStateChannel fsm.StateChannels, SyncChan nc.NetworkChannels) {
	var (
		elevList        [config.NumElevator]config.Elev 
		OnlineList      [config.NumElevator]bool        
		TempKeyOrder    config.Keypress                 
		TempButtonEvent elevio.ButtonEvent             
	)

	for {
		select {
		case newLocalOrder := <-HardwareToControl: //new order from hardware
			id := costFunction(Local_ID, newLocalOrder, elevList, OnlineList)
			if id != -1 {
				if id == Local_ID {
					LocalStateChannel.NewOrder <- newLocalOrder //send order local
				} else {
					TempKeyOrder = config.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button}
					SyncChan.LocalOrderToExternal <- TempKeyOrder // send orders abroad
				}
			}
		case TempKeyOrder = <-SyncChan.ExternalOrderToLocal: //Receive order from external
			TempButtonEvent = elevio.ButtonEvent{Button: TempKeyOrder.Button, Floor: TempKeyOrder.Floor}
			if elevList[Local_ID].State == config.Undefined {
				costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)
				TempKeyOrder.DesignatedElevator = costID
				SyncChan.LocalOrderToExternal <- TempKeyOrder
			} else {
				LocalStateChannel.NewOrder <- TempButtonEvent
			}
		case NewUpdateLocalElevator := <-LocalStateChannel.Elevator: //Receive a change in local elevator from fsm
			change := false
			for floor := 0; floor < config.NumFloor; floor++ {
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					if elevList[Local_ID].Queue[floor][button] && !NewUpdateLocalElevator.Queue[floor][button] { //if local elevator had finished an order
						change = true 
					}
				}
				if change {
					change = false
					for id := 0; id < config.NumElevator; id++ {
						if id != Local_ID {
							elevList[id].Queue[floor][elevio.BT_HallUp] = false
							elevList[id].Queue[floor][elevio.BT_HallDown] = false
						}
					}
				}
			}
			change = false
			if elevList[Local_ID].State != config.Undefined && NewUpdateLocalElevator.State == config.Undefined { //if local elevator experienses engine failure
				elevList[Local_ID].State = config.Undefined
				for floor := 0; floor < config.NumFloor; floor++ {
					for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
						if NewUpdateLocalElevator.Queue[floor][button] { 
							TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}                       
							costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)                  
							TempKeyOrder = config.Keypress{Floor: floor, Button: button, DesignatedElevator: costID} 
							SyncChan.LocalOrderToExternal <- TempKeyOrder                                            
						}
					}
				}
			}
			elevList[Local_ID] = NewUpdateLocalElevator
			go func() { UpdateLight <- elevList }()
			if OnlineList[Local_ID] {
				go func() { SyncChan.LocalElevatorToExternal <- elevList }()
			}
		case tempElevatorArray := <-SyncChan.UpdateMainLogic: //Receive change of external elevators from network controller
			change := false
			tempQueue := elevList[Local_ID].Queue
			for id := 0; id < config.NumElevator; id++ {
				if id != Local_ID {
					for floor := 0; floor < config.NumFloor; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if elevList[id].Queue[floor][button] && !tempElevatorArray[id].Queue[floor][button] { //if external elevator had finished an order
								change = true 
							}
						}
						if change {
							change = false
							for newID := 0; newID < config.NumElevator; newID++ {
								if newID == Local_ID {
									tempQueue[floor][elevio.BT_HallUp] = false
									tempQueue[floor][elevio.BT_HallDown] = false
								}
								if newID != id && newID != Local_ID {
									tempElevatorArray[newID].Queue[floor][elevio.BT_HallUp] = false
									tempElevatorArray[newID].Queue[floor][elevio.BT_HallDown] = false
								}
							}
						}
					}
				}
			}
			if tempQueue != elevList[Local_ID].Queue {
				elevList[Local_ID].Queue = tempQueue
				go func() { LocalStateChannel.DeleteQueue <- elevList[Local_ID].Queue }()
				if OnlineList[Local_ID] {
					go func() { SyncChan.LocalElevatorToExternal <- elevList }()
				}
			}
			for id := 0; id < config.NumElevator; id++ {
				if id == Local_ID {
					continue
				}
				elevList[id] = tempElevatorArray[id]
			}
			go func() { UpdateLight <- elevList }()
		case NewOnlineList := <-SyncChan.OnlineElevators: //Receive change in who's on the networks from network controller
			change := false
			numofOnlineElevs := 0

			for id := 0; id < config.NumElevator; id++ {
				if NewOnlineList[id] == true {
					numofOnlineElevs++
				}
			}
			if numofOnlineElevs == 0 { //if local elevator is offline
				for id := 0; id < config.NumElevator; id++ {
					for floor := 0; floor < config.NumFloor; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if id != Local_ID {
								change = true
								elevList[id].Queue[floor][button] = false
							}
						}
						elevList[Local_ID].Queue[floor][elevio.BT_HallUp] = false
						elevList[Local_ID].Queue[floor][elevio.BT_HallDown] = false
					}
				}
				if change {
					LocalStateChannel.DeleteQueue <- elevList[Local_ID].Queue
					go func() { SyncChan.LocalElevatorToExternal <- elevList }()
				}
			}
			change = false
			if numofOnlineElevs > 0 { // if external elevator is offline
				for id := 0; id < config.NumElevator; id++ {
					if OnlineList[id] && !NewOnlineList[id] { 
						for floor := 0; floor < config.NumFloor; floor++ {
							for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
								if elevList[id].Queue[floor][button] {
									change = true
									elevList[id].Queue[floor][button] = false
									if button != elevio.BT_Cab {
										TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
										costID := costFunction(Local_ID, TempButtonEvent, elevList, NewOnlineList)
										if costID == Local_ID {
											LocalStateChannel.NewOrder <- TempButtonEvent
										}
									}
								}
							}
						}
					}
				}
			}
			if change {
				go func() { SyncChan.LocalElevatorToExternal <- elevList }()
			}
			go func() { UpdateLight <- elevList }()
			OnlineList = NewOnlineList
		}
	}
}

func LightSetter(Local_ID int, UpdateLight chan [config.NumElevator]config.Elev) { //update Hall and Cab order lights
	var Order [config.NumElevator]bool
	for {
		select {
		case Elevator := <-UpdateLight:
			for floor := 0; floor < config.NumFloor; floor++ {
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					for id := 0; id < config.NumElevator; id++ {
						Order[id] = false
						if id != Local_ID && button == elevio.BT_Cab {
							// Ignore inside orders for other elevators
							continue
						}
						if Elevator[id].Queue[floor][button] {
							elevio.SetButtonLamp(button, floor, true)
							Order[id] = true
						}
					}
					if Order == [config.NumElevator]bool{false} {
						elevio.SetButtonLamp(button, floor, false)
					}
				}
			}
		}
	}
}
