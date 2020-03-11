package main

import (
	"./config"
	"./elevio"
	"./fsm"
)

func main() {

	var numFloors int = 4

	elevio.Init("localhost:15657", numFloors)
	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	receiveOrders := make(chan elevio.ButtonEvent)
	//receiveFloors := make(chan int)
	channels := fsm.StateChannels{
		OrderComplete:  make(chan int),
		Elevator:       make(chan config.Elev),
		NewOrder:       make(chan config.Keypress),
		ArrivedAtFloor: make(chan int),
	}

	go elevio.PollButtons(receiveOrders)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)

	go fsm.UpdateKeys(channels.NewOrder, receiveOrders)

	go fsm.RunElevator(channels)
	//go OH.UpdateHallAndCabButtons(receiveOrders)
	//go OH.UpdateFloor(receiveFloors)

	for {
		//fmt.Println(timer.CheckTime(t))
	}

}
