package main

import (
	//"fmt"
	"strconv"
	//"time"
	"os"

	"./config"
	"./elevio"
	"./fsm"
	ec "./ElevatorController"
	nc"./NetworkController"
	bcast "./network/bcast"
	peers "./network/peers"
)


func main() {
	// For å kjøre koden må du først skrive "run go main.go" som vanlig så må du spesifisere IDen til heisen og hvilke port den skal kobles til
	// på simulatoren
	//eksempel kode: go run main.go 0 15000
	//koden kjører med id = 0 og port = 15000
	 
	LocalIDString := os.Args[1]
	localhost := "localhost:" + os.Args[2]
	LocalID,err := strconv.Atoi(LocalIDString)

	if err != nil {
		panic(err.Error())
	}
	
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
	msgpPort := 42030 //Har bare valgt en random port for å teste kode.
	orderPort := 42050

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

	for{ 
	}
}

	
