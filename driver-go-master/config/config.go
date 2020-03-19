package config

import (
	"../elevio"
)

const (
	NumFloor    int = 4
	NumElevator int = 3
	NumButtons  int = 3
)

/*type ButtonType int AKA Button

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)*/

/*type MotorDirection int  AKA Direction

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)*/

type Acknowledge int

const (
	Finished Acknowledge = iota - 1
	NotAck
	Acked
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
	Button                elevio.ButtonType
	DesignatedElevator int
	Completed          bool
}

type Elev struct {
	State ElevState
	Dir   elevio.MotorDirection
	Floor int
	Queue [NumFloor][NumButtons]bool
}

type AckList struct {
	ElevatorID		   int
	ImplicitAcks       [NumElevator]Acknowledge
}

type Message struct {
	Elevator         [NumElevator]Elev
	RegisteredOrders [NumFloor][NumButtons - 1]AckList
	ID               int
}
