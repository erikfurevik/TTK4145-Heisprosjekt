package NetworkController

import (
	"../config"
	//"../elevio"
	//"../fsm"
	peers "../network/peers"
	//bcast "../network/bcast"
	"strconv"
	"time"
	"fmt"
)


type NetworkChannels struct {
	//from network to elevator controller
	UpdateMainLogic  			chan [config.NumElevator]config.Elev //updates elevator controller function with the states of the other elevators
	OnlineElevators 			chan [config.NumElevator]bool	//updats elevator controller channel used to send the status of the online elevators from sync to gov
	ExternalOrderToLocal		chan config.Keypress			//updates elevator controller.Order we get from abroad that we pass on to elevator controller
	
	//from elevator to network controller
	LocalOrderToExternal  		chan config.Keypress 			//channel used to send orders to other elevators
	//LocalElevatorToExternal  	chan config.Elev 				//channel that send the status of the local elevator from gov to sync
	LocalElevatorToExternal		chan [config.NumElevator]config.Elev

	//network controller to network
	OutgoingMsg     			chan config.Message			//not cocern of gov
	OutgoingOrder 				chan config.Keypress		// new order from elevator controller going to network through the network controller
	PeersTransmitEnable    		chan bool					//channel going to network, updating the other elevators about my presence

	//network to network controller
	IncomingMsg     			chan config.Message			//not concern of gov
	IncomingOrder 				chan config.Keypress		//NewOrder from abroad to network controller
	PeerUpdate     	 			chan peers.PeerUpdate		//channel going to network controller update about the other networks
	
}

func NetworkController(Local_ID int, channel NetworkChannels){
	var (
		//registeredOrders [config.NumFloor][config.NumElevator - 1]config.Acklist
		//elevList 		[config.NumElevator]config.Elev
		msg 			config.Message
		onlineList 		[config.NumElevator]bool
		//recentlyDies 	[config.NumElevator]bool
		//someUpdate 		bool
		//offline 		bool
	)
	


	//lostID := -1
	reassignTimer := time.NewTimer(5 *time.Second)
	broadcastElevTimer  	:= time.NewTimer(100 * time.Millisecond)
	singleModeTicker 	:= time.NewTicker(100 *time.Millisecond)


	reassignTimer.Stop()
	broadcastElevTimer.Stop()
	singleModeTicker.Stop()
	

	//In the future, i will have to get acknoledged for everything i send from all elevators
	//I will have to acknolede for everything that i receive. 
	//I will have to check whos on the network and  
	//I will have to update the other elevator about my presence
	//Only send new data if the previous have been been acknowledged
	//onlineList = [config.NumElevator]bool {true}
	msg.ID = Local_ID
	channel.PeersTransmitEnable <- true



	for {
		select {
		case msg.Elevator = <- channel.LocalElevatorToExternal: //update of our elevator
			if onlineList[Local_ID]{
				channel.OutgoingMsg <- msg
					//fmt.Println("send local elevator state: ", Local_ID)
			}
		case ExternalOrder := <- channel.LocalOrderToExternal: //get order from controller
			channel.OutgoingOrder <- ExternalOrder //send it over the network
			//fmt.Println("send local order to abroad")

		case inOrder := <- channel.IncomingOrder: //order from network
		if inOrder.DesignatedElevator == Local_ID {
			channel.ExternalOrderToLocal <- inOrder
			//channel.ExternalOrderToLocal <- inOrder
			//fmt.Println("receive local order")
		}
		case inMSG := <- channel.IncomingMsg: //state of an elevator abroad
			fmt.Println(inMSG.ID)
			if inMSG.ID != Local_ID &&  inMSG != msg{
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID] //update message strcut
				channel.UpdateMainLogic <- msg.Elevator
				//fmt.Println("receive external elevator:", inMSG.ID)
			}
		case peerUpdate := <-channel.PeerUpdate:

			if len(peerUpdate.Peers) == 0{
				for id := 0; id < config.NumElevator; id++ {
					onlineList[id] = false				
				}
			}
			if len(peerUpdate.New) > 0 {
				newElev, _ := strconv.Atoi(peerUpdate.New)
				onlineList[newElev] = true
			}
			if len(peerUpdate.Lost) > 0 {
				lostElev, _ := strconv.Atoi(peerUpdate.Lost[0])
				onlineList[lostElev] = false
			}

			go func () {channel.OnlineElevators <- onlineList} ()

		
			fmt.Println("Number peers.", len(peerUpdate.Peers))
			fmt.Println("New peers: ", peerUpdate.New)
			fmt.Println("Lost peers", peerUpdate.Lost)
		}
	}
}