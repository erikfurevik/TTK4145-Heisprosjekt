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

	//orderCleared := false
	//channel.Elevator <- elevator

	for {
		select {
		case newOrder := <-channel.NewOrder:
			fmt.Println("New order")
			/*
				if newOrder.Completed {
					elevator.Queue[newOrder.Floor][elevio.BT_HallUp] = false
					elevator.Queue[newOrder.Floor][elevio.BT_HallDown] = false
					orderCleared = true

				} else {
					elevator.Queue[newOrder.Floor][newOrder.Button] = true
				}
			*/
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
					//channel.OrderComplete <- newOrder.Floor
					//go func() { channel.OrderComplete <- newOrder.Floor }()

				} else {
					elevator.State = config.Moving
					EngineFailureTimer.Reset(3 * time.Second)
				}
			//case config.Moving:
			case config.DoorOpen:
				if elevator.Floor == newOrder.Floor {
					DoorTimer.Reset(3 * time.Second)
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{false}
					//channel.OrderComplete <- newOrder.Floor
				}
			case config.Undefined:
				fmt.Println("fatal error")
			}
			//channel.Elevator <- elevator

		case elevator.Floor = <-channel.ArrivedAtFloor:
			if shouldMotorStop(elevator) {
				//orderCleared = false
				elevio.SetDoorOpenLamp(true)
				EngineFailureTimer.Stop()
				elevator.State = config.DoorOpen
				elevio.SetMotorDirection(elevio.MD_Stop)
				DoorTimer.Reset(3 * time.Second)
				elevator.Queue[elevator.Floor] = [config.NumButtons]bool{false}
				//channel.OrderComplete <- elevator.Floor

			} else if elevator.State == config.Moving {
				EngineFailureTimer.Reset(3 * time.Second)
			}
			//channel.Elevator <- elevator

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
			//channel.Elevator <- elevator
		case <-EngineFailureTimer.C:
			//elevio.SetMotorDirection(elevio.MD_Stop)
			elevator.State = config.Undefined
			fmt.Println("Engine failure")
			//elevio.SetMotorDirection(elevator.Dir)
			//channel.Elevator <- elevator
			EngineFailureTimer.Reset(5 * time.Second)

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
