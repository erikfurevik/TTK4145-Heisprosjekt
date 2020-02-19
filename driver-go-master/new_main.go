package main

import (
	"fmt"
	"time"

	"./elevio"
)

func main() {

	type MotorDirection int

	const (
		MD_Up   MotorDirection = 1
		MD_Down                = -1
		MD_Stop                = 0
	)

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	// var currentDirection elevio.MotorDirection = elevio.MD_Stop
	// Use elevio.SetMotorDirection(dir)

	DrvButtons := make(chan elevio.ButtonEvent)
	DrvFloors := make(chan int)

	go elevio.PollButtons(DrvButtons)    // Polls request buttons
	go elevio.PollFloorSensor(DrvFloors) // Checks current elev floor

	var nextRequest int = -1

	type ElevatorState int
	const (
		ELEV_Idle   ElevatorState = 1
		ELEV_Open                 = 2
		ELEV_Moving               = 3
	)

	currentFloor := -1

	go func() { // FSM?

		currentState := ELEV_Idle
		inState := false

		for {
			// GET NEXT REQUEST FROM QUEUE: nextRequest =
			switch currentState {
			case ELEV_Idle:
				if !inState {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(false)
					inState = true
				}
				if nextRequest != -1 {
					if nextRequest == currentFloor {
						currentState = ELEV_Open
					}
					if nextRequest != currentFloor {
						currentState = ELEV_Moving
					}
					inState = false
				}
			case ELEV_Open:
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				nextRequest = -1
				elevio.SetButtonLamp(elevio.BT_HallUp, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_HallDown, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_Cab, currentFloor, false)
				time.Sleep(time.Second * 3)
				elevio.SetDoorOpenLamp(false)
				currentState = ELEV_Idle
			case ELEV_Moving:
				if !inState {
					if nextRequest < currentFloor {
						elevio.SetMotorDirection(elevio.MD_Down)
						currentState = ELEV_Moving
					}
					if nextRequest > currentFloor {
						elevio.SetMotorDirection(elevio.MD_Up)
						currentState = ELEV_Moving
					}
					inState = true
				}
				if currentFloor == nextRequest {
					currentState = ELEV_Open
					inState = false
					elevio.SetMotorDirection(elevio.MD_Stop)
				}
			} // End switch
		} // End for
	}() // End FSM

	for {
		select {
		case a := <-DrvButtons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			// SHOULD ADD TO QUEUE, MEANWHILE:
			nextRequest = a.Floor

		case a := <-DrvFloors:
			fmt.Printf("%+v\n", a)
			elevio.SetFloorIndicator(a)
			currentFloor = a
			// Update floor

		} // End select
	} // End for
} // End main
