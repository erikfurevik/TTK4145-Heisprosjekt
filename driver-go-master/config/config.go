package config

import (
	"../elevio"
)

const (
	NumFloor    int = 9
	NumElevator int = 3
	NumButtons  int = 3
)

type ElevState int

const (
	Undefined ElevState = iota - 1
	Idle
	Moving
	DoorOpen
)
type Keypress struct {
	Floor              int
	Button             elevio.ButtonType
	DesignatedElevator int
}
type Elev struct {
	State ElevState
	Dir   elevio.MotorDirection
	Floor int
	Queue [NumFloor][NumButtons]bool
}


type Message struct {
	Elevator         [NumElevator]Elev
	ID               int
}
