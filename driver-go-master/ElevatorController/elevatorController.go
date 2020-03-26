package ElevatorController

import (
	"../config"
	"../elevio"
	"../fsm"
	nc "../NetworkController"
	"fmt"
)

//type StateChannels struct {
//	OrderComplete chan int
//	ArrivedAtFloor chan int
//	NewOrder chan elevio.ButtonEvent
//	Elevator chan config.Elev
//	DeleteNewOrder chan elevio.ButtonEvent
//}

//type NetworkChannels struct {
//	UpdateMainLogic  			chan [config.NumElevator]config.Elev //updates the governor function with the states of the other elevators
//	LocalElevatorToExternal  	chan config.Elev 				//channel that send the status of the local elevator from gov to sync
//	LocalOrderToExternal  		chan config.Keypress //channel used to send orders to other elevators
//	ExternalOrderToLocal		chan config.Keypress
//	OnlineElevators 			chan [config.NumElevator]bool	//channel used to send the status of the online elevators from sync to gov
//	IncomingMsg     			chan config.Message			//not concern of gov
//	OutgoingMsg     			chan config.Message			//not cocern of gov
//	//PeerUpdate      chan peers.PeerUpdate	//not concern of gov
//	//PeerTxEnable    chan bool				//not concern of gov
//}



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
		OnlineList = [config.NumElevator]bool {true}
		fmt.Println("starting mainlogic function:", Local_ID)


		for {
			select {
			case newLocalOrder := <- HardwareToControl:
				fmt.Println("new order from hardware:",Local_ID)
				//there is a new order from the hardware
				if !OnlineList[Local_ID] { //our elevator is ofline -> take only cab
					if newLocalOrder.Button == elevio.BT_Cab{
						LocalStateChannel.NewOrder <- newLocalOrder //send order local
					}
				}else if OnlineList[Local_ID] && elevList[Local_ID].State == config.Undefined{
					//We are online but motor is not working
					if newLocalOrder.Button != elevio.BT_Cab{
						tempOnline := OnlineList
						tempOnline[Local_ID] = false
						id := costFunction(Local_ID, newLocalOrder, elevList, tempOnline)
						TempKeyOrder = config.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button, Completed: false}
						SyncChan.LocalOrderToExternal <- TempKeyOrder // send orders abroad
					}else{
						LocalStateChannel.NewOrder <- newLocalOrder //send order local
					}

				}else if OnlineList[Local_ID] && elevList[Local_ID].State != config.Undefined{ 
					//local elevator is working normally
					if newLocalOrder.Floor == elevList[Local_ID].Floor && elevList[Local_ID].State != config.Moving {
						LocalStateChannel.NewOrder <- newLocalOrder //send order local
					}else{
						id := costFunction(Local_ID, newLocalOrder, elevList, OnlineList)
						if id == Local_ID {
							LocalStateChannel.NewOrder <- newLocalOrder //send order local
						}else {
							TempKeyOrder = config.Keypress{DesignatedElevator: id, Floor: newLocalOrder.Floor, Button: newLocalOrder.Button, Completed: false}
							SyncChan.LocalOrderToExternal <- TempKeyOrder // send orders abroad
						}
					}
				}


			case TempKeyOrder = <- SyncChan.ExternalOrderToLocal:
				TempButtonEvent= elevio.ButtonEvent{Button: TempKeyOrder.Button, Floor: TempKeyOrder.Floor}
				//if our own engine is down, recalculate cost for order and send back
				if elevList[Local_ID].State == config.Undefined{
					costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)
					TempKeyOrder.DesignatedElevator = costID
					go func() {SyncChan.LocalOrderToExternal <- TempKeyOrder}()
				}else {
					go func() {LocalStateChannel.NewOrder <- TempButtonEvent} () //send order local
				
				}


			case NewUpdateLocalElevator := <- LocalStateChannel.Elevator:
				//Update about our new local elevator
				if elevList[Local_ID].State != config.Undefined && NewUpdateLocalElevator.State == config.Undefined{ //i am undefiend now
					for floor := 0; floor < config.NumFloor; floor++{
						for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
							if elevList[Local_ID].Queue[floor][button]{
								TempButtonEvent= elevio.ButtonEvent{Floor: floor, Button: button}
								costID := costFunction(Local_ID, TempButtonEvent, elevList, OnlineList)
								TempKeyOrder = config.Keypress{Floor: floor, Button: button, DesignatedElevator: costID}
								elevList[Local_ID].Queue[floor][button] = false
								SyncChan.LocalOrderToExternal <- TempKeyOrder //send order external
							}
						}
					}
				}
				NewUpdateLocalElevator.Queue = elevList[Local_ID].Queue
				elevList[Local_ID] = NewUpdateLocalElevator //update info about elevator
				UpdateLight <- elevList//update lights
				if OnlineList[Local_ID] {
					SyncChan.LocalElevatorToExternal <- elevList[Local_ID]
				}


			case TempKeyOrder.Floor = <- LocalStateChannel.OrderComplete: 
				//an order is complete from the local fsm
				fmt.Println("order completed:",Local_ID)
				//Delete the finisehd Hall Button orders from all elevators
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					if elevList[Local_ID].Queue[TempKeyOrder.Floor][button] {
						TempKeyOrder.Button = button
					}
					for elevator := 0; elevator < config.NumElevator; elevator++ {
						if button != elevio.BT_Cab || elevator == Local_ID {
							elevList[elevator].Queue[TempKeyOrder.Floor][button] = false
						}
					}
				}
				UpdateLight <- elevList //update lights
				
				if OnlineList[Local_ID]{
					SyncChan.LocalElevatorToExternal <- elevList[Local_ID]
				}


			case tempElevatorArray := <-SyncChan.UpdateMainLogic:
				//The update from the other elevators about their state, queue, motor direction etc
				change := false
				//If previousl Hall Order is completed by another elevator, Delete that order from our fsm as well
				for floor := 0; floor < config.NumFloor; floor++{
					for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
						if elevList[Local_ID].Queue[floor][button]{
							for id := 0; id < config.NumElevator; id++ {
								if id != Local_ID {
									if elevList[id].Queue[floor][button] && !tempElevatorArray[id].Queue[floor][button]{
										elevList[Local_ID].Queue[floor][button] = false
										order := elevio.ButtonEvent{Floor: floor, Button: button}
										LocalStateChannel.DeleteNewOrder <- order
										change = true
									}
								}
							}
						}
					}
				}
				for id := 0; id < config.NumElevator; id++ {
					if id == Local_ID {
						continue
					}
					if elevList[id].Queue != tempElevatorArray[id].Queue {
						change = true
					}
					elevList[id] = tempElevatorArray[id] //See if there are any chagnes. Save the updated elevators
				}
				if change {
					UpdateLight <- elevList
				}

			case NewOnlineList := <- SyncChan.OnlineElevators:
				//if another elevator goes offline
				for id := 0; id < config.NumElevator; id++{
					if id != Local_ID && OnlineList[id] && !NewOnlineList[id]{
						for floor := 0; floor < config.NumFloor; floor++ {
							for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
								if elevList[id].Queue[floor][button]{
									TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
									costID := costFunction(Local_ID, TempButtonEvent, elevList, NewOnlineList)
									elevList[id].Queue[floor][button] = false
									if costID == Local_ID{
										LocalStateChannel.NewOrder <- TempButtonEvent //might have to make a go func out of this, might have to youse a buffer
										elevList[Local_ID].Queue[floor][button] = true //Might not be necessary,
									}
								}
							}
						}
					}
				}
				//if our elevator goes offline
				if OnlineList[Local_ID] && !NewOnlineList[Local_ID]{
					for floor := 0; floor < config.NumFloor; floor++ {
						for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
							if elevList[Local_ID].Queue[floor][button]{
								TempButtonEvent = elevio.ButtonEvent{Floor: floor, Button: button}
								LocalStateChannel.DeleteNewOrder <- TempButtonEvent
							}
						}
					}
				}
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