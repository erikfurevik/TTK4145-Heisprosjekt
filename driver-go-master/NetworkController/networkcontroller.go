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
	LocalElevatorToExternal		chan [config.NumElevator]config.Elev //channel that send the status of the local elevator from gov to sync

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
		msg 			config.Message
		onlineList 		[config.NumElevator]bool
	)
	
	reassignTimer 			:= time.NewTimer(5 *time.Second)
	broadcastMsgTicker  	:= time.NewTicker(60 * time.Millisecond)

	reassignTimer.Stop()
	//broadcastMsgTicker.Reset(100 * time.Millisecond)
	

	//In the future, i will have to get acknoledged for everything i send from all elevators
	//I will have to acknolede for everything that i receive. 
	//Only send new data if the previous have been been acknowledged
	msg.ID = Local_ID
	channel.PeersTransmitEnable <- true


	for {
		select {
		case msg.Elevator = <- channel.LocalElevatorToExternal: //update of our elevator
			
		case ExternalOrder := <- channel.LocalOrderToExternal: //get order from controller
			channel.OutgoingOrder <- ExternalOrder //send it over the network

		case inOrder := <- channel.IncomingOrder: //order from network
		if inOrder.DesignatedElevator == Local_ID {
			channel.ExternalOrderToLocal <- inOrder
		}
		case inMSG := <- channel.IncomingMsg: //state of an elevator abroad
			if inMSG.ID != Local_ID &&  inMSG.Elevator[inMSG.ID] != msg.Elevator[inMSG.ID]{
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID] //update message strcut
				channel.UpdateMainLogic <- msg.Elevator
				fmt.Println("Receiving")
			}
		case <- broadcastMsgTicker.C:
			if onlineList[Local_ID]{
				channel.OutgoingMsg <- msg
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