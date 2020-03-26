package main

import (
	//"fmt"
	//"strconv"
	//"time"

	"./config"
	"./elevio"
	"./fsm"
	ec "./ElevatorController"
	nc"./NetworkController"
	bcast "./network/bcast"
	peers "./network/peers"
)


func runParallel(LocalID int, localhost string){
	//defualt := "localhost:15657"

	elevio.Init(localhost, config.NumFloor)


	channels := fsm.StateChannels{
		OrderComplete:  make(chan int),
		Elevator:       make(chan config.Elev),
		NewOrder:       make(chan elevio.ButtonEvent),
		ArrivedAtFloor: make(chan int),
	}

	network := nc.NetworkChannels{
		//from network to elevator controller
	UpdateMainLogic:  			make(chan [config.NumElevator]config.Elev),
	OnlineElevators: 			make(chan [config.NumElevator]bool),	
	ExternalOrderToLocal:		make(chan config.Keypress),			
	
	//from elevator to network controller
	LocalOrderToExternal:  		make(chan config.Keypress), 			
	LocalElevatorToExternal:  	make(chan config.Elev), 			
	
	//network controller to network
	OutgoingMsg:     			make(chan config.Message),		
	OutgoingOrder: 				make(chan config.Keypress),		
	PeerTxEnable:    			make(chan bool),				

	//network to network controller
	IncomingMsg:     			make(chan config.Message),			
	IncomingOrder: 				make(chan config.Keypress),		
	PeerUpdate:     	 		make(chan peers.PeerUpdate),
	} 

	var (
		newOrder = make(chan elevio.ButtonEvent)
		updateLight = make(chan [config.NumElevator]config.Elev)
	)


	//id_string := strconv.Itoa(LocalID)
	msgpPort := 42034 //Har bare valgt en random port for å teste kode.
	orderPort := 42035

	go elevio.PollButtons(newOrder)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)
	
	go fsm.RunElevator(channels)
	go ec.MainLogicFunction(LocalID ,newOrder, updateLight, channels, network)
	go ec.LightSetter(updateLight,LocalID)
	go nc.NetworkController(LocalID, network)


	go bcast.Transmitter(msgpPort, network.OutgoingMsg)           
	go bcast.Receiver(msgpPort, network.IncomingMsg)
	
	go bcast.Transmitter(orderPort, network.OutgoingOrder)           
	go bcast.Receiver(orderPort, network.IncomingOrder)


	//go peers.Transmitter(port, id_string, enableTx) 
	//go peers.Receiver(port, peerUpdate)             

}


func main() {
	go runParallel(0, "localhost:20000")
	go runParallel(1, "localhost:15000")

	for {

	}

}
