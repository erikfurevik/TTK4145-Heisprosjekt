package controller

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


			case Floor := <- LocalStateChannel.OrderComplete: 
				//an order is complete from the local fsm
				elevList[Local_ID].Queue[Floor] = [3]bool{false}
				go func() {UpdateLight <- elevList} ()//update lights


			case tempElevatorArray := <-SyncChan.UpdateMainLogic:
				//The update from the other elevators about their state, queue, motor direction etc
			

			case OnlineList = <- SyncChan.OnlineElevators:
				//The online states of the elevators

			}
		}
	}