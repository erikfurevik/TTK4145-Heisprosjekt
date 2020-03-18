package ElevatorController

import (
	"../config"
	"../elevio"
	"../fsm"
)
type SyncChannels struct {
	UpdateMainLogic  			chan [config.NumElevator]config.Elev //updates the governor function with the states of the other elevators
	LocalElevatorToExternal  	chan config.Elev 				//channel that send the status of the local elevator from gov to sync
	LocalOrderToExternal  		chan config.Keypress //channel used to send orders to other elevators
	ExternalOrderToLocal		chan config.Keypress
	OnlineElevators 			chan [config.NumElevator]bool	//channel used to send the status of the online elevators from sync to gov
	//IncomingMsg     chan config.Message			//not concern of gov
	//OutgoingMsg     chan config.Message			//not cocern of gov
	//PeerUpdate      chan peers.PeerUpdate	//not concern of gov
	//PeerTxEnable    chan bool				//not concern of gov
}

//type StateChannels struct {
//	OrderComplete chan int
//	ArrivedAtFloor chan int
//	NewOrder chan elevio.ButtonEvent
//	Elevator chan config.Elev
//	DeleteNewOrder chan elevio.ButtonEvent
//}


/*takes in status of the other online elevators, takes in the status of local elevator, takes in new orders from hardarw
  and decides what should be done */
func MainLogicFunction(Local_ID int, HardwareToControl <-chan elevio.ButtonEvent, UpdateLight chan<- [config.NumElevator]config.Elev,
	LocalStateChannel fsm.StateChannels, SyncChan SyncChannels){
		var (
			elevList 		[config.NumElevator]config.Elev //Our Total info about all elevators
			OnlineList 		[config.NumElevator]bool		//check which elevator that is online
			
			completeOrder 	config.Keypress					//dont know what this is yet
			
			TempKeyOrder 	config.Keypress					//Helper struct to Convert between ButtonEvent and Keypress
			TempButtonEventOrder elevio.ButtonEvent			//Helper struct to convert between ButonEvent and Keypress

		)
		//completeOrder.DesignatedElevator = Local_ID;
		elevList[Local_ID] = <- LocalStateChannel.Elevator //Blockin until we get info aboout our elevator
		SyncChan.LocalElevatorToExternal <- elevList[Local_ID] //Blocking until external elevator has recevied our elevator info

		for {
			select {
			case newLocalOrder := <- HardwareToControl: 
				//there is a new order from the hardware
				if !OnlineList[Local_ID] { //our elevator is ofline -> Single elevator mode
					//elevList[Local_ID].Queue[newLocalOrder.Floor][newLocalOrder.Button] = true //update queue
					//go func() {UpdateLight <- elevList} ()//update lights
					go func() {LocalStateChannel.NewOrder <- newLocalOrder} () //send order local
				}else if OnlineList[Local_ID] && elevList[Local_ID].State == config.Undefined{
					//We are online but motor is not working
					//if not cab button send order external, if cab send internal
					if newLocalOrder.Button != elevio.BT_Cab{
						//not Cab
						tempOnline := OnlineList
						tempOnline[Local_ID] = false
						id := costFunction(Local_ID, newLocalOrder, elevList, tempOnline)
						TempKeyOrder.DesignatedElevator = id	//ID of the elevator that should take the order
						TempKeyOrder.Floor = newLocalOrder.Floor
						TempKeyOrder.Button = newLocalOrder.Button
						//go func() {UpdateLight <- elevList} ()//update lights
						go func() {SyncChan.LocalOrderToExternal <- TempKeyOrder} () // send orders abroad
					}else{
						go func() {LocalStateChannel.NewOrder <- newLocalOrder} () //send order local
					}
				}else if OnlineList[Local_ID] && elevList[Local_ID].State != config.Undefined{ //local elevator is working normally
					if newLocalOrder.Floor == elevList[Local_ID].Floor && elevList[Local_ID].State != config.Moving {
						//elevList[Local_ID].Queue[newLocalOrder.Floor][newLocalOrder.Button] = true //update Quee matrix
						//go func() {UpdateLight <- elevList} ()//update lights
						go func() {LocalStateChannel.NewOrder <- newLocalOrder} () //send order local
					}else{
						id := costFunction(Local_ID, newLocalOrder, elevList, OnlineList)
						if id == Local_ID {
							//elevList[Local_ID].Queue[newLocalOrder.Floor][newLocalOrder.Button] = true
							//go func() {UpdateLight <- elevList} () //update lights
							go func() {LocalStateChannel.NewOrder <- newLocalOrder} () //send order local
						}else {
							TempKeyOrder.DesignatedElevator = id	//ID of the elevator that should take the order
							TempKeyOrder.Floor = newLocalOrder.Floor
							TempKeyOrder.Button = newLocalOrder.Button
							//go func() {UpdateLight <- elevList} ()//update lights
							go func() {SyncChan.LocalOrderToExternal <- TempKeyOrder} () // send orders abroad
						}
					}
				}
			case newExternalOrder := <- SyncChan.ExternalOrderToLocal:
				// we receive an order from abroad
				//elevList[Local_ID].Queue[newExternalOrder.Floor][newExternalOrder.Button] = true
				TempButtonEventOrder.Button = newExternalOrder.Button
				TempButtonEventOrder.Floor = newExternalOrder.Floor
				go func() {LocalStateChannel.NewOrder <- TempButtonEventOrder} () //send order local

				
			case NewUpdateLocalElevator := <- LocalStateChannel.Elevator:
				//Update from the fsm of the local elevator
				//temp = elevList[Local_ID].Queue
				elevList[Local_ID] = NewUpdateLocalElevator //update info about elevator
				//elevList[Local_ID].Queue = elevList[Local_ID].
				go func() {UpdateLight <- elevList} ()//update lights
				if OnlineList[Local_ID] {
					SyncChan.LocalElevatorToExternal <- elevList[Local_ID]
				}


			case completeOrder.Floor = <- LocalStateChannel.OrderComplete: 
				//an order is complete from the local fsm
				//update the the other elevators as well

				//Delete the orderer finished from all elevtors
				for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++ {
					if elevList[Local_ID].Queue[completeOrder.Floor][button] {
						completeOrder.Button = button
					}
					for elevator := 0; elevator < config.NumElevator; elevator++ {
						if button != elevio.BT_Cab || elevator == Local_ID {
							elevList[elevator].Queue[completeOrder.Floor][button] = false
						}
					}
				}

				if OnlineList[Local_ID]{
					SyncChan.LocalElevatorToExternal <- elevList[Local_ID]
				}
				go func() {UpdateLight <- elevList} () //update lights


			case tempElevatorArray := <-SyncChan.UpdateMainLogic:
				//The update from the other elevators about their state, queue, motor direction etc
				change := false
				//The part we update external elevators
				for elevator := 0; elevator < config.NumElevator; elevator++ {
					if elevator == Local_ID {
						continue
					}
					if elevList[elevator].Queue != tempElevatorArray[elevator].Queue {
						change = true
					}
					elevList[elevator] = tempElevatorArray[elevator] //Se if there are any chagnes. Save the updated elevators
				}
				//The part we update our own elevator
				for floor := 0; floor < config.NumFloor; floor++{
					for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
						//Our local didnt have an order before, but it does now
						if tempElevatorArray[Local_ID].Queue[floor][button] && !elevList[Local_ID].Queue[floor][button] { 
							//then we update the Q matrix for our local elevator
							elevList[Local_ID].Queue[floor][button] = true
							order := elevio.ButtonEvent{Floor: floor, Button: button}
							go func() { LocalStateChannel.NewOrder <- order }()
							//If our local had order previously but it does not have it now
							change = true
						}else if !tempElevatorArray[Local_ID].Queue[floor][button] && elevList[Local_ID].Queue[floor][button]{
							elevList[Local_ID].Queue[floor][button] = false
							order := elevio.ButtonEvent{Floor: floor, Button: button}
							go func() { LocalStateChannel.DeleteNewOrder <- order }()
							change = true
						}
					}

				}
				if change {
					UpdateLight <- elevList
				}


			case OnlineList = <- SyncChan.OnlineElevators:
				//The online states of the elevators

			}
		}
	}