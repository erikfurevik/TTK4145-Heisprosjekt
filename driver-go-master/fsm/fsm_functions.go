package fsm

import (
	"../config"
	"../elevio"
)

func orderAbove(elevator Elev) bool {
	for floor := elevator.Floor + 1; floor < config.NumFloor; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func orderBelow(elevator Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func shouldMotorStop(elevator Elev) {
	switch elevator.Dir {
	case elevio.BT_Up:
		return elevator.Queue[elevator.Floor][elevio.BT_HallUp] ||
			elevator.Queue[elevator.Floor][elevio.BT_Cab] ||
			!orderAbove(elevator)
	case elevio.BT_Down:
		return elevator.Queue[elevator.Floor][elevio.BT_HallDown] ||
			elevator.Queue[elevator.Floor][elevio.BT_Cab] ||
			!orderBelow(elevator)
	default:
	}
	return false
}

func chooseDirection(elevator elev) elevio.MotorDirection {
	switch elevator.Dir {
	case elevio.MD_Stop:
		if orderAbove(elevator) {
			return elevio.MD_Up
		} else if orderBelow(elevator) {
			return elevio.MD_Down
		} else {
			elevio.MD_Stop
		}
	case elevio.MD_Up:
		if orderAbove(elevator) {
			return elevio.MD_Up
		} else if orderBelow(elevator) {
			return elevio.MD_Down
		} else {
			elevio.MD_Stop
		}
	case elevio.MD_Down:
		if orderBelow(elevator) {
			return elevio.MD_Down
		} else if orderAbove(elevator) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}
