package main

import (
	"./elevio"
	"./fsm"
)

func main() {

	var numFloors int = 4

	elevio.Init("localhost:15657", numFloors)
	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	receiveOrders := make(chan elevio.ButtonEvent)
	receiveFloors := make(chan int)
	channels := fsm.StateChannels{}

	go elevio.PollButtons(receiveOrders)
	go elevio.PollFloorSensor(receiveFloors)

	go fsm.UpdateKeys(channels, receiveOrders, receiveFloors)

	go fsm.RunElevator(channels)
	//go OH.UpdateHallAndCabButtons(receiveOrders)
	//go OH.UpdateFloor(receiveFloors)

	for {
		//fmt.Println(timer.CheckTime(t))
	}

}
