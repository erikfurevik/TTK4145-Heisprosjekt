package fsm

import (
	"fmt"
	"time"

	"../config"
	"../elevio"
)

type StateChannels struct {
	OrderComplete chan int
	ArrivedAtFloor chan int
	NewOrder chan elevio.ButtonEvent
	Elevator chan config.Elev
	DeleteNewOrder chan elevio.ButtonEvent
	DeleteQueue chan [config.NumFloor][config.NumButtons] bool 
}

func RunElevator(channel StateChannels) {
	elevator := config.Elev{
		State: config.Idle,
		Dir:   elevio.MD_Stop,
		Floor: elevio.GetFloor(),
	}
	DoorTimer := time.NewTimer(3 * time.Second)
	EngineFailureTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()
	EngineFailureTimer.Stop()
	updateExternal := false

	//var newOrder config.ButtonEvent
	//TakeOrderTimer := := time.NewTimer(3 * time.Second)

	//channel.Elevator <- elevator

	for {
		select {
		case newOrder := <-channel.NewOrder:
			elevator.Queue[newOrder.Floor][newOrder.Button] = true
			switch elevator.State {
			case config.Idle:
				elevator.Dir = chooseDirection(elevator)
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir == elevio.MD_Stop {
					elevator.State = config.DoorOpen
					elevio.SetDoorOpenLamp(true)
					DoorTimer.Reset(3 * time.Second)
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{false}

				} else {
					elevator.State = config.Moving
					EngineFailureTimer.Reset(3 * time.Second)
				}
				updateExternal = true
			
			case config.Moving:
				updateExternal = true

			case config.DoorOpen:
				if elevator.Floor == newOrder.Floor {
					DoorTimer.Reset(3 * time.Second)
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{false}
				}else{
					updateExternal = true
				}
				
			case config.Undefined:
				fmt.Println("fatal error")
				updateExternal = true
			}
			
		case deleteOrder := <- channel.DeleteNewOrder:
			elevator.Queue[deleteOrder.Floor][deleteOrder.Button] = false 

		case deleteQueue := <- channel.DeleteQueue:
			elevator.Queue = deleteQueue

		case elevator.Floor = <-channel.ArrivedAtFloor:
			if shouldMotorStop(elevator) {
				EngineFailureTimer.Stop()
				elevio.SetMotorDirection(elevio.MD_Stop)
				if !orderAtFloor(elevator){
					elevator.State = config.Idle
					DoorTimer.Reset(3 * time.Millisecond)
				}else {
					elevio.SetDoorOpenLamp(true)
					elevator.State = config.DoorOpen
					DoorTimer.Reset(3 * time.Second)
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{false}

				}
			} else if elevator.State == config.Moving {
				EngineFailureTimer.Reset(3 * time.Second)
			}
			updateExternal = true
		case <-DoorTimer.C:
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = chooseDirection(elevator)
			if elevator.Dir == elevio.MD_Stop {
				elevator.State = config.Idle
				EngineFailureTimer.Stop()
			} else {
				elevator.State = config.Moving
				EngineFailureTimer.Reset(3 * time.Second)
				elevio.SetMotorDirection(elevator.Dir)
			}
			//channel.OrderComplete <- elevator.Floor
			updateExternal = true
		case <-EngineFailureTimer.C:
			elevator.State = config.Undefined
			fmt.Println("Engine failure")
			EngineFailureTimer.Reset(5 * time.Second)
			updateExternal = true
		}
		if updateExternal{
			channel.Elevator <- elevator
			updateExternal = false
		}
	}
}

//UpdateKeys ..
func UpdateKeys(NewOrder chan config.Keypress, receiveOrder chan elevio.ButtonEvent) {
	var key config.Keypress
	key.DesignatedElevator = 1
	key.Completed = false
	for {
		select {
		case order := <-receiveOrder:
			key.Floor = order.Floor
			key.Button = order.Button
			//fmt.Println(key.Floor)
			NewOrder <- key

		}
	}
}

//Testchannels ..
func Testchannels(channel StateChannels) {
	for {
		select {
		case a := <-channel.ArrivedAtFloor:
			fmt.Println(a)
		case b := <-channel.NewOrder:
			fmt.Println(b.Floor)
		}
	}
}
