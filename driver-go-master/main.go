package main

import (
	"./elevio"
)

func main() {

	var numFloors int = 4

	elevio.Init("localhost:15657", numFloors)
	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	/*
		receiveOrders := make(chan elevio.ButtonEvent)
		receiveFloors := make(chan int)
		drv_obstr := make(chan bool)
		drv_stop := make(chan bool)

		go elevio.PollButtons(receiveOrders)
		go elevio.PollFloorSensor(receiveFloors)
		go elevio.PollObstructionSwitch(drv_obstr)
		go elevio.PollStopButton(drv_stop)

		go OH.UpdateHallAndCabButtons(receiveOrders)
		go OH.UpdateFloor(receiveFloors)
	*/
	//fsm.TEST()

	elevio.SetMotorDirection(elevio.MD_Up)

}
