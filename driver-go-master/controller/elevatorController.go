package controller

import (
	"../config"
	"../elevio"
	"../fsm"
)


func Governate(Local_ID int, HardwareToControl chan config.Keypress, UpdateLight chan [config.NumElevator],
	channel fsm.StateChannels, SyncChannel chan Keypress)