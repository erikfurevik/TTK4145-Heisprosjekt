package fsm

import (
	"fmt"
	"time"

	"../config"
	"../elevio"
)

type StateChannels struct {
	ArrivedAtFloor chan int
	NewOrder       chan elevio.ButtonEvent
	Elevator       chan config.Elev
	DeleteQueue    chan [config.NumFloor][config.NumButtons]bool
}

/*This function handles the fsm logic of the elevator. It has the following respisibilities
	- Take orders in the correct sequence
	- Update elevator controller about any changes in local elevator
	- Check for engine failure
	- Handle fault tolerance such as system crash by saving orders to file
	- Handle floor lights
*/
func RunElevator(LocalID int,channel StateChannels) {
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

	for elevio.GetFloor() == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	readFromFile("cabOrders", LocalID, &elevator)

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
				} else {
					updateExternal = true
				}
			case config.Undefined:
				fmt.Println("fatal error")
				updateExternal = true
			}
		case deleteQueue := <-channel.DeleteQueue:
			elevator.Queue = deleteQueue
		case elevator.Floor = <-channel.ArrivedAtFloor:
			elevio.SetFloorIndicator(elevator.Floor)
			if shouldMotorStop(elevator) {
				EngineFailureTimer.Stop()
				elevio.SetMotorDirection(elevio.MD_Stop)
				if !orderAtFloor(elevator) {
					elevator.State = config.Idle
					DoorTimer.Reset(3 * time.Millisecond)
				} else {
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
			updateExternal = true
		case <-EngineFailureTimer.C:
			elevator.State = config.Undefined
			fmt.Println("Engine failure")
			EngineFailureTimer.Reset(5 * time.Second)
			updateExternal = true
		}
		if updateExternal {
			updateExternal = false
			writetoFile("cabOrders", LocalID, elevator)
			go func() { channel.Elevator <- elevator }()
		}
	}
}
