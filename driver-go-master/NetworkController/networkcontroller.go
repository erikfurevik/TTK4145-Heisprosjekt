package NetworkController

import (
	"../config"
	//"../elevio"
	//"../fsm"
)


type NetworkChannels struct {
	UpdateMainLogic  			chan [config.NumElevator]config.Elev //updates the governor function with the states of the other elevators
	
	LocalElevatorToExternal  	chan config.Elev 				//channel that send the status of the local elevator from gov to sync
	LocalOrderToExternal  		chan config.Keypress //channel used to send orders to other elevators
	ExternalOrderToLocal		chan config.Keypress
	OnlineElevators 			chan [config.NumElevator]bool	//channel used to send the status of the online elevators from sync to gov
	IncomingMsg     			chan config.Message			//not concern of gov
	OutgoingMsg     			chan config.Message			//not cocern of gov
	//PeerUpdate      chan peers.PeerUpdate	//not concern of gov
	//PeerTxEnable    chan bool				//not concern of gov
}

