package fsm

import (
	"../config"
	"../elevio"
)

func orderAbove(elevator config.Elev) bool {
	for floor := elevator.Floor + 1; floor < config.NumFloor; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func orderBelow(elevator config.Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func shouldMotorStop(elevator config.Elev) bool {
	switch elevator.Dir {
	case elevio.MD_Up:
		return elevator.Queue[elevator.Floor][elevio.BT_HallUp] ||
			elevator.Queue[elevator.Floor][elevio.BT_Cab] ||
			!orderAbove(elevator)
	case elevio.MD_Down:
		return elevator.Queue[elevator.Floor][elevio.BT_HallDown] ||
			elevator.Queue[elevator.Floor][elevio.BT_Cab] ||
			!orderBelow(elevator)
	case elevio.MD_Stop:
	default:
	}
	return false
}

func chooseDirection(elevator config.Elev) elevio.MotorDirection {
	switch elevator.Dir {
	case elevio.MD_Stop:
		if orderAbove(elevator) {
			return elevio.MD_Up
		} else if orderBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Up:
		if orderAbove(elevator) {
			return elevio.MD_Up
		} else if orderBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
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
