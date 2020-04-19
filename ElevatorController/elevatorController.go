package ElevatorController

import (
	"fmt"

	nc "../NetworkController"
	"../config"
	"../elevio"
	"../fsm"
)

/*takes in status of the other online elevators, takes in the status of local elevator, takes in new orders from hardarware
  and decides what should be done */
func MainLogicFunction(Local_ID int, HardwareToControl <-chan elevio.ButtonEvent, UpdateLight chan<- [config.NumElevator]config.Elev,
	LocalStateChannel fsm.StateChannels, SyncChan nc.NetworkChannels) {
	var (
		elevList        [config.NumElevator]config.Elev //total info about all elevators
		OnlineList      [config.NumElevator]bool        //online elevators
		TempKeyOrder    config.Keypress                 //helper struct to Convert between ButtonEvent and Keypress
		TempButtonEvent elevio.ButtonEvent              //helper struct to convert between ButtonEvent and Keypress
	)
	fmt.Println("starting mainlogic function:", Local_ID)

	for {
		select {
		case newLocalOrder := <-HardwareToControl:
			id := costFunction(Local_ID, newLocalOrder, elevList, OnlineList)
			if id != -1 {
				if id == Local_ID {
					LocalStateChannel.NewOrder <- newLocalOrder //send order local
				} else {
					TempKeyOrder = config.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button}
					SyncChan.LocalOrderToExternal <- TempKeyOrder // send orders abroad
				}
			}
		case TempKeyOrder = <-SyncChan.ExternalOrderToLocal:
			TempButtonEvent = elevio.ButtonEvent{Button: TempKeyOrder.Button, Floor: TempKeyOrder.Floor}
			//if our engine is down, recalculate cost for order and send back
			if elevList[Local_ID].State == config.Undefined {
				costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)
				TempKeyOrder.DesignatedElevator = costID
				SyncChan.LocalOrderToExternal <- TempKeyOrder
			} else {
				LocalStateChannel.NewOrder <- TempButtonEvent //send order local
			}
		case NewUpdateLocalElevator := <-LocalStateChannel.Elevator:
			//Update about our new local elevator
			change := false
			for floor := 0; floor < config.NumFloor; floor++ {
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					if elevList[Local_ID].Queue[floor][button] && !NewUpdateLocalElevator.Queue[floor][button] {
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
			if elevList[Local_ID].State != config.Undefined && NewUpdateLocalElevator.State == config.Undefined {
				elevList[Local_ID].State = config.Undefined
				for floor := 0; floor < config.NumFloor; floor++ {
					for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ { //distribute all orders except cab
						if NewUpdateLocalElevator.Queue[floor][button] { //i have hall order at that floor
							TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}                       //make ButtonEvent
							costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)                  //calculate cost
							TempKeyOrder = config.Keypress{Floor: floor, Button: button, DesignatedElevator: costID} //buttonEvent into keypress
							SyncChan.LocalOrderToExternal <- TempKeyOrder                                            //broadcast order
						}
					}
				}
			}
			elevList[Local_ID] = NewUpdateLocalElevator //update info about elevator
			go func() { UpdateLight <- elevList }()
			if OnlineList[Local_ID] {
				go func() { SyncChan.LocalElevatorToExternal <- elevList }()
			}
		case tempElevatorArray := <-SyncChan.UpdateMainLogic:
			change := false
			tempQueue := elevList[Local_ID].Queue
			for id := 0; id < config.NumElevator; id++ {
				if id != Local_ID {
					for floor := 0; floor < config.NumFloor; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if elevList[id].Queue[floor][button] && !tempElevatorArray[id].Queue[floor][button] {
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
				elevList[id] = tempElevatorArray[id] //save updated info about elevators
			}
			go func() { UpdateLight <- elevList }()
		case NewOnlineList := <-SyncChan.OnlineElevators:
			change := false
			numofOnlineElevs := 0

			for id := 0; id < config.NumElevator; id++ {
				if NewOnlineList[id] == true {
					numofOnlineElevs++
				}
			}
			if numofOnlineElevs == 0 { //elevator is offline
				for id := 0; id < config.NumElevator; id++ {
					for floor := 0; floor < config.NumFloor; floor++ {
						for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
							if id != Local_ID {
								change = true
								elevList[id].Queue[floor][button] = false //delete all order from all elevators
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
			if numofOnlineElevs > 0 { //another elevator is offline
				for id := 0; id < config.NumElevator; id++ {
					if OnlineList[id] && !NewOnlineList[id] { // for every elevator that if offline
						for floor := 0; floor < config.NumFloor; floor++ {
							for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
								if elevList[id].Queue[floor][button] { //for every hall order that they had
									change = true
									elevList[id].Queue[floor][button] = false
									if button != elevio.BT_Cab {
										TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
										costID := costFunction(Local_ID, TempButtonEvent, elevList, NewOnlineList)
										if costID == Local_ID { //recalculate the cost
											LocalStateChannel.NewOrder <- TempButtonEvent //orders local elevator should take
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

func LightSetter(UpdateLight chan [config.NumElevator]config.Elev, Local_ID int) {
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
