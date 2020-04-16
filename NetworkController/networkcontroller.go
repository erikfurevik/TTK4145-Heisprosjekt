package NetworkController

import (
	"fmt"
	"strconv"
	"time"

	"../config"
	peers "../network/peers"
)

type NetworkChannels struct {
	//from network to elevator controller
	UpdateMainLogic      chan [config.NumElevator]config.Elev //updates elevatorController with the states of the other elevators.
	OnlineElevators      chan [config.NumElevator]bool        //updated elevatorController of which elevators are online
	ExternalOrderToLocal chan config.Keypress                 //updates elevator controller with orders recieved from abroad.

	//from elevator to network controller
	LocalOrderToExternal    chan config.Keypress                 //channel used to broadcast local orders.
	LocalElevatorToExternal chan [config.NumElevator]config.Elev //channel that broadcast the status of the local elevator

	//network controller to network
	OutgoingMsg         chan config.Message  //not cocern of elevatorController
	OutgoingOrder       chan config.Keypress //new local order that is broadcasted
	PeersTransmitEnable chan bool            //channel updating the other elevators about my presence

	//network to network controller
	IncomingMsg   chan config.Message   //not concern of elevatorController
	IncomingOrder chan config.Keypress  //new order recieved from abroad elevator
	PeerUpdate    chan peers.PeerUpdate //channel going to network controller update about the other networks

}

func NetworkController(Local_ID int, channel NetworkChannels) {
	var (
		msg           config.Message
		onlineList    [config.NumElevator]bool
		outgoingOrder config.Keypress
		incomingOrder config.Keypress
	)

	PrimaryOrderTimer := time.NewTicker(100 * time.Millisecond)
	orderTicker := time.NewTicker(10 * time.Millisecond)
	broadcastMsgTicker := time.NewTicker(40 * time.Millisecond)
	deleteIncomingOrderTicker := time.NewTicker(1 * time.Second)
	orderTicker.Stop()

	msg.ID = Local_ID
	channel.PeersTransmitEnable <- true
	queue := make([]config.Keypress, 0)

	for {
		select {
		case msg.Elevator = <-channel.LocalElevatorToExternal: //update of our elevator
		case ExternalOrder := <-channel.LocalOrderToExternal: //recived order on local elevator
			queue = append(queue, ExternalOrder) //put order into queue

		case inOrder := <-channel.IncomingOrder: //recieved order from abroad elevator
			if inOrder.DesignatedElevator == Local_ID && incomingOrder != inOrder {
				incomingOrder = inOrder
				channel.ExternalOrderToLocal <- inOrder
			}
		case inMSG := <-channel.IncomingMsg: //state of an abroad elevator
			if inMSG.ID != Local_ID && inMSG.Elevator[inMSG.ID] != msg.Elevator[inMSG.ID] {
				msg.Elevator[inMSG.ID] = inMSG.Elevator[inMSG.ID] //update message strcut
				channel.UpdateMainLogic <- msg.Elevator
			}
		case <-broadcastMsgTicker.C:
			if onlineList[Local_ID] {
				channel.OutgoingMsg <- msg
			}
		case <-PrimaryOrderTimer.C:
			if len(queue) > 0 {
				outgoingOrder = queue[0]
				queue = queue[1:]
				orderTicker = time.NewTicker(10 * time.Millisecond)
			} else {
				orderTicker.Stop()
			}
		case <-orderTicker.C:
			channel.OutgoingOrder <- outgoingOrder //bradcast order

		case <-deleteIncomingOrderTicker.C:
			incomingOrder = config.Keypress{Floor: -1}
		case peerUpdate := <-channel.PeerUpdate:
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

			fmt.Println("Number peers.", len(peerUpdate.Peers))
			fmt.Println("New peers: ", peerUpdate.New)
			fmt.Println("Lost peers", peerUpdate.Lost)
		}
	}
}
