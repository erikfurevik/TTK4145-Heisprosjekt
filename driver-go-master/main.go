package main

import (
	"./config"
	"./elevio"
	"./fsm"
)

func main() {

	elevio.Init("localhost:15657", config.NumFloor)
	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	//receiveOrders := make(chan elevio.ButtonEvent)
	//receiveFloors := make(chan int)
	channels := fsm.StateChannels{
		OrderComplete:  make(chan int),
		Elevator:       make(chan config.Elev),
		NewOrder:       make(chan elevio.ButtonEvent),
		ArrivedAtFloor: make(chan int),
	}

	go elevio.PollButtons(channels.NewOrder)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)

	go fsm.RunElevator(channels)

	for {
		//fmt.Println(timer.CheckTime(t))
	}

}
