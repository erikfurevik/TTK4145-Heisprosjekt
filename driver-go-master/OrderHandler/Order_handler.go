package OrderHandler

import (
	"fmt"

	"../config"
	"../elevio"
)

/*const (
	BtnUP
	BtnDown
	BtnInside
)*/

type ArrayID struct {
	HallButtons int
	CabButons   int
	Floor       int
	ID          int
	State       int
}

var ID [config.NumElevator]int

//var State [config.NumElevator]int
var Floor [config.NumElevator]int

var HallButtonsElevator1 [config.NumFloor][2]int

var CabButtons [config.NumFloor][config.NumFloor]int

/*this function continously updates HallButtons and CabButtons*/
func UpdateHallAndCabButtons(receiver chan<- elevio.ButtonEvent) {
	for {
		select {
		case data := <-receiver:
			if data.Button == elevio.BT_Cab {
				CabButtons[data.Floor][0] = 1
			}
			if data.Button == elevio.BT_HallUp {
				HallButtonsElevator1[data.Floor][1] = 1 //trur jeg
			}
			if data.Button == elevio.BT_HallDown {
				HallButtonsElevator1[data.Floor][0] = 1 //trur jeg
			}
		}
		fmt.Println(HallButtonsElevator1)
	}
}

func UpdateFloor(receiver chan int) {
	for {
		select {
		case data := <-receiver:
			Floor[0] = data + 1
			fmt.Println(Floor)
		}
	}
}

func ReturnFloor(col int) int {
	return Floor[col]
}

func ReturnHallButtons(row int, col int) int {
	return HallButtons[row][col]
}

func ReturnCabButtons(row int, col int) int {
	return CabButtons[row][col]
}

func WriteHalButtons(row int, col int, value int) {
	HallButtons[row][col] = value
}

func WriteCabButtons(row int, col int, value int) {
	CabButtons[row][col] = value
}
