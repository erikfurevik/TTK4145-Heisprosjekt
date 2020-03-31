package ElevatorController

import (
	"../config"
	"../elevio"
	"../fsm"
	nc "../NetworkController"
	"fmt"
)



/*takes in status of the other online elevators, takes in the status of local elevator, takes in new orders from hardarw
  and decides what should be done */
func MainLogicFunction(Local_ID int, HardwareToControl <-chan elevio.ButtonEvent, UpdateLight chan<- [config.NumElevator]config.Elev,
	LocalStateChannel fsm.StateChannels, SyncChan nc.NetworkChannels){
		var (
			elevList 		[config.NumElevator]config.Elev //Our Total info about all elevators
			OnlineList 		[config.NumElevator]bool		//check which elevator that is online
			TempKeyOrder 	config.Keypress					//Helper struct to Convert between ButtonEvent and Keypress
			TempButtonEvent elevio.ButtonEvent			//Helper struct to convert between ButonEvent and Keypress
		)
		fmt.Println("starting mainlogic function:", Local_ID)



		for {
			select {
			case newLocalOrder := <- HardwareToControl:
				id := costFunction(Local_ID, newLocalOrder, elevList, OnlineList)
				if id != -1 {
					if id == Local_ID {
						LocalStateChannel.NewOrder <- newLocalOrder //send order local
					}else {
						TempKeyOrder = config.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button}
						SyncChan.LocalOrderToExternal <- TempKeyOrder // send orders abroad
					}
				}
			case TempKeyOrder = <- SyncChan.ExternalOrderToLocal:
				TempButtonEvent= elevio.ButtonEvent{Button: TempKeyOrder.Button, Floor: TempKeyOrder.Floor}
				//if our own engine is down, recalculate cost for order and send back
				if elevList[Local_ID].State == config.Undefined{
					costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)
					TempKeyOrder.DesignatedElevator = costID
					SyncChan.LocalOrderToExternal <- TempKeyOrder
				}else {
					LocalStateChannel.NewOrder <- TempButtonEvent //send order local
				}
			case NewUpdateLocalElevator := <- LocalStateChannel.Elevator:
				//Update about our new local elevator
				change := false
				for floor := 0; floor < config.NumFloor; floor++{
					for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
						if elevList[Local_ID].Queue[floor][button] && !NewUpdateLocalElevator.Queue[floor][button]{
							change = true
						}
					}
					if change {
						change = false
						for id := 0; id < config.NumElevator; id++{
							if id != Local_ID{
								elevList[id].Queue[floor][elevio.BT_HallUp] = false
								elevList[id].Queue[floor][elevio.BT_HallDown] = false
							}	
						}
					}
				}				
				change = false
				//i will have to test this further with more elevators
				if elevList[Local_ID].State != config.Undefined && NewUpdateLocalElevator.State == config.Undefined{ //i am undefiend now 
					elevList[Local_ID].State = config.Undefined //update that we are undefined
					for floor := 0; floor < config.NumFloor; floor++{
						for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{//distribute all orders except cab
							if NewUpdateLocalElevator.Queue[floor][button]{ //i have hall order at that floor
								TempButtonEvent= elevio.ButtonEvent{Floor: floor, Button: button} //make Button struct
								costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList) //calculate cost
								TempKeyOrder = config.Keypress{Floor: floor, Button: button, DesignatedElevator: costID} //but into keypress
								//elevList[costID].Queue[floor][button] = true //at into their queue, maybe no necessary
								//NewUpdateLocalElevator.Queue[floor][button] = false
								SyncChan.LocalOrderToExternal <- TempKeyOrder //send order external
							}
						}
					}
				}
				elevList[Local_ID] = NewUpdateLocalElevator //update info about elevator
				go func () {UpdateLight <- elevList} ()
				if OnlineList[Local_ID] {
					go func() {SyncChan.LocalElevatorToExternal <- elevList} ()
				}
			case tempElevatorArray := <-SyncChan.UpdateMainLogic:
				change := false
				tempQueue := elevList[Local_ID].Queue
				for id := 0; id < config.NumElevator; id++ {
					if id != Local_ID {
						for floor := 0; floor < config.NumFloor; floor++{
							for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
								if elevList[id].Queue[floor][button] && !tempElevatorArray[id].Queue[floor][button] {
									change = true
								}
							}
							if change{
								change = false
								for newID := 0; newID < config.NumElevator; newID++{
									if newID == Local_ID{
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
					go func() {LocalStateChannel.DeleteQueue <- elevList[Local_ID].Queue} ()
					if OnlineList[Local_ID] {
					go func() {SyncChan.LocalElevatorToExternal <- elevList} ()
					}
				}
				for id := 0; id < config.NumElevator; id++ {
					if id == Local_ID {
						continue
					}
					elevList[id] = tempElevatorArray[id] //See if there are any chagnes. Save the updated elevators
				}
				go func () {UpdateLight <- elevList} ()
			case NewOnlineList := <- SyncChan.OnlineElevators:
				change := false
				numofOnlineElevs := 0

				for id := 0; id < config.NumElevator; id++ {
					if NewOnlineList[id] == true {
						numofOnlineElevs  ++
					}				
				}
				if numofOnlineElevs == 0{ // we are offline
					for id := 0; id < config.NumElevator; id++ {
						for floor := 0; floor < config.NumFloor; floor++ {
							for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
								if id != Local_ID {
									change = true
									elevList[id].Queue[floor][button] = false //delete all order from all elevators
								}
							}
							elevList[Local_ID].Queue[floor][elevio.BT_HallUp] = false 
							elevList[Local_ID].Queue[floor][elevio.BT_HallDown]	= false
						}
					}
					if change{
						LocalStateChannel.DeleteQueue <-elevList[Local_ID].Queue
						go func() {SyncChan.LocalElevatorToExternal <- elevList} ()
					}	
				}
				change = false
				if numofOnlineElevs > 0{ //i am online but someone else fell off
					for id := 0; id < config.NumElevator; id++{ 
						if OnlineList[id] && !NewOnlineList[id]{ // for every elevator that fell off
							for floor := 0; floor < config.NumFloor; floor++ {
								for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
									if elevList[id].Queue[floor][button]{ //for every hall order that they had
										change = true
										elevList[id].Queue[floor][button] = false
										if button != elevio.BT_Cab {
											TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
											costID := costFunction(Local_ID, TempButtonEvent, elevList, NewOnlineList)
											//elevList[id].Queue[floor][button] = false //delete the order from their queue
											if costID == Local_ID{ //recalculate the cost and if we should take the orders, then we send it to ourselves. 
												LocalStateChannel.NewOrder <- TempButtonEvent 
												//elevList[Local_ID].Queue[floor][button] = true //Might not be necessary,
											}
										}
									}
								}															
							}
						}
					}
				}
				if change {
					go func() {SyncChan.LocalElevatorToExternal <- elevList} ()
				}
				go func () {UpdateLight <- elevList} ()
				OnlineList = NewOnlineList
			}
		}
	}


func LightSetter(UpdateLight chan [config.NumElevator]config.Elev, Local_ID int) {
	var Order [config.NumElevator]bool
	for {
		select{
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