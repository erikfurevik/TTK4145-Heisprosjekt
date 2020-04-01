package fsm

import (
	"../config"
	"../elevio"
	"os"
	"fmt"
	"strconv"
)

func orderAbove(elevator config.Elev) bool {
	for floor := elevator.Floor + 1; floor < config.NumFloor; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
			if elevator.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func orderBelow(elevator config.Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
			if elevator.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}
func orderAtFloor(elevator config.Elev)bool{
	for button := elevio.BT_HallUp; button <= elevio.BT_Cab; button++{
		if elevator.Queue[elevator.Floor][button]{
			return true
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
		return true
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


func writetoFile(filname string, LocalID int, elevator config.Elev){
	idstring := strconv.Itoa(LocalID)
	writeFile, _ := os.Create(filname + idstring)
	var stringVariable string
	
	for i := 0; i < config.NumFloor; i++{
		if elevator.Queue[i][elevio.BT_Cab] {
			stringVariable = stringVariable + "1"
		}else{
			stringVariable = stringVariable + "0"
			}	
	}
	fmt.Println(stringVariable)
	data := []byte(stringVariable)
	writeFile.Write(data)

}


func readFromFile(filename string ,LocalID int, elevator * config.Elev){
	idstring := strconv.Itoa(LocalID)
	readFile, _ := os.Open(filename + idstring)
	data := make([]byte, config.NumFloor)
	readFile.Read(data)

	for i := 0; i < config.NumFloor; i++{
		if string(data[i]) == "1" {
			elevator.Queue[i][elevio.BT_Cab] = true
		}else{
			elevator.Queue[i][elevio.BT_Cab] = false
		}
	}
	readFile.Close()

	fmt.Println(elevator.Queue)
}
