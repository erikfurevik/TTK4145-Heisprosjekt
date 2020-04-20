package NetworkController

import (
	"strconv"
	"time"

	"../config"
	peers "../network/peers"
)

type NetworkChannels struct {
	//from network controller to elevator controller
	UpdateMainLogic      chan [config.NumElevator]config.Elev 	//updates elevatorController with the states of the other elevators.
	OnlineElevators      chan [config.NumElevator]bool        	//updated elevatorController of which elevators are online
	ExternalOrderToLocal chan config.Keypress                 	//updates elevator controller with orders recieved from abroad.

	//from elevator controller to network controller
	LocalOrderToExternal    chan config.Keypress                 //broadcast local orders.
	LocalElevatorToExternal chan [config.NumElevator]config.Elev //broadcast the status of the local elevator

	//network controller to network
	OutgoingMsg         chan config.Message  					
	OutgoingOrder       chan config.Keypress 					
	PeersTransmitEnable chan bool            					

	//network to network controller
	IncomingMsg   chan config.Message   						
	IncomingOrder chan config.Keypress  						
	PeerUpdate    chan peers.PeerUpdate 						

}
/*This is the brain of the network. It has the following resposibilities
	- Check which elevator is on the network
	- Pass local orders to external elevator
	- Pass external orders to local elevator
	- Receive external elevator updates
	- Send local elevator updates
	- It handles fault tolerances such as packet loss and packet corruption
*/

func NetworkController(Local_ID int, channel NetworkChannels) {
	var (
		msg           config.Message
		onlineList    [config.NumElevator]bool
		outgoingOrder config.Keypress
		incomingOrder config.Keypress
	)

	primaryOrderTicker := time.NewTicker(100 * time.Millisecond)
	orderTicker := time.NewTicker(10 * time.Millisecond)
	broadcastMsgTicker := time.NewTicker(40 * time.Millisecond)
	deleteIncomingOrderTicker := time.NewTicker(1 * time.Second)
	orderTicker.Stop()

	msg.ID = Local_ID
	channel.PeersTransmitEnable <- true
	queue := make([]config.Keypress, 0)

	for {
		select {
		case msg.Elevator = <-channel.LocalElevatorToExternal: //update of our elevator struct

		case ExternalOrder := <-channel.LocalOrderToExternal: //recived order on local elevator
			queue = append(queue, ExternalOrder)

		case inOrder := <-channel.IncomingOrder: //receive order from abroad
			if inOrder.DesignatedElevator == Local_ID && incomingOrder != inOrder {
				incomingOrder = inOrder
				channel.ExternalOrderToLocal <- inOrder
			}
		case inMSG := <-channel.IncomingMsg: //incomming abroad elevator struct
			if inMSG.ID != Local_ID && inMSG.Elevator[inMSG.ID] != msg.Elevator[inMSG.ID] {
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID]
				channel.UpdateMainLogic <- msg.Elevator
			}
		case <-broadcastMsgTicker.C: //broadcast our elevator struct
			if onlineList[Local_ID] {
				channel.OutgoingMsg <- msg
			}
		case <-primaryOrderTicker.C: // reset order ticker if queue is not empty, else stop order ticker
			if len(queue) > 0 {
				outgoingOrder = queue[0]
				queue = queue[1:]
				orderTicker = time.NewTicker(10 * time.Millisecond)
			} else {
				orderTicker.Stop()
			}
		case <-orderTicker.C: //send the same order for every tick
			channel.OutgoingOrder <- outgoingOrder

		case <-deleteIncomingOrderTicker.C: //clean the incomingOrder variable at every tick
			incomingOrder = config.Keypress{Floor: -1}

		case peerUpdate := <-channel.PeerUpdate: //check online status
			if len(peerUpdate.Peers) == 0 {
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
			go func() { channel.OnlineElevators <- onlineList }()
		}
	}
}
