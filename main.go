package main

import (
	"os"
	"strconv"

	ec "./ElevatorController"
	nc "./NetworkController"
	"./config"
	"./elevio"
	"./fsm"
	bcast "./network/bcast"
	peers "./network/peers"
)

func main() {

	LocalIDString := os.Args[1]
	localhost := "localhost:" + os.Args[2]
	LocalID, _ := strconv.Atoi(LocalIDString)

	msgpPort  := 42000
	orderPort := 43000
	peersPort := 44000

	elevio.Init(localhost, config.NumFloor)
	channels := fsm.StateChannels{
		Elevator:       			make(chan config.Elev),
		NewOrder:       			make(chan elevio.ButtonEvent, 100),
		ArrivedAtFloor: 			make(chan int),
		DeleteQueue:    			make(chan [config.NumFloor][config.NumButtons]bool),
	}

	network := nc.NetworkChannels{
		UpdateMainLogic:      		make(chan [config.NumElevator]config.Elev, 100),
		OnlineElevators:      		make(chan [config.NumElevator]bool),
		ExternalOrderToLocal: 		make(chan config.Keypress),

		LocalOrderToExternal:    	make(chan config.Keypress),
		LocalElevatorToExternal: 	make(chan [config.NumElevator]config.Elev),

		OutgoingMsg:         		make(chan config.Message),
		OutgoingOrder:       		make(chan config.Keypress),
		PeersTransmitEnable: 		make(chan bool),

		IncomingMsg:   				make(chan config.Message, 30),
		IncomingOrder: 				make(chan config.Keypress),
		PeerUpdate:    				make(chan peers.PeerUpdate),
	}

	var (
		newOrder    = make(chan elevio.ButtonEvent)
		updateLight = make(chan [config.NumElevator]config.Elev)
	)
	
	go elevio.PollButtons(newOrder)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)

	go fsm.RunElevator(LocalID, channels)
	go ec.MainLogicFunction(LocalID, newOrder, updateLight, channels, network)
	go nc.NetworkController(LocalID, network)
	go ec.LightSetter(LocalID, updateLight)

	go bcast.Transmitter(msgpPort, network.OutgoingMsg)
	go bcast.Receiver(msgpPort, network.IncomingMsg)

	go bcast.Transmitter(orderPort, network.OutgoingOrder)
	go bcast.Receiver(orderPort, network.IncomingOrder)

	go peers.Transmitter(peersPort, LocalIDString, network.PeersTransmitEnable)
	go peers.Receiver(peersPort, network.PeerUpdate)

	select {}
}
