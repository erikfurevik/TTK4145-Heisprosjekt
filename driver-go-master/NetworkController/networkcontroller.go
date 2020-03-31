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
		outgoingOrder 	config.Keypress
		incomingOrder  	config.Keypress
	)
	
	PrimaryOrderTimer 			:= time.NewTicker(100 * time.Millisecond)
	orderTicker					:= time.NewTicker(10 * time.Millisecond)
	broadcastMsgTicker  		:= time.NewTicker(40 * time.Millisecond)
	deleteIncomingOrderTicker 	:= time.NewTicker(1 * time.Second)
	orderTicker.Stop()

	msg.ID = Local_ID
	channel.PeersTransmitEnable <- true
	queue := make([]config.Keypress, 0) 


	for {
		select {
		case msg.Elevator = <- channel.LocalElevatorToExternal: //update of our elevator
		case ExternalOrder := <- channel.LocalOrderToExternal: //get order from controller
			queue = append(queue, ExternalOrder) //put order to queue

		case inOrder := <- channel.IncomingOrder: //order from network
		if inOrder.DesignatedElevator == Local_ID && incomingOrder != inOrder {
			incomingOrder = inOrder
			channel.ExternalOrderToLocal <- inOrder
			fmt.Println("incoming order")
		}
		case inMSG := <- channel.IncomingMsg: //state of an elevator abroad
			if inMSG.ID != Local_ID &&  inMSG.Elevator[inMSG.ID] != msg.Elevator[inMSG.ID]{
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID] //update message strcut
				channel.UpdateMainLogic <- msg.Elevator
			}
		case <- broadcastMsgTicker.C:
			if onlineList[Local_ID]{
				channel.OutgoingMsg <- msg
			}
		case <- PrimaryOrderTimer.C:
			if len(queue) > 0{
				outgoingOrder = queue[0];
				queue = queue[1:]
				orderTicker = time.NewTicker(10 * time.Millisecond)
			}else {
				orderTicker.Stop()
			}
		case <- orderTicker.C:
			channel.OutgoingOrder <- outgoingOrder //send it over the network
			fmt.Println("sending order")

		case <- deleteIncomingOrderTicker.C:
			incomingOrder = config.Keypress {Floor: -1,}
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