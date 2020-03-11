package fsm

import (
	"fmt"
	"time"

	"../config"
	"../elevio"
)

const floors int = 4

/*value 1 = up
value 2 = down
value 3 = cab
*/

/*så fort kost funksjonen bergner om en ordre skal som skal bli tatt lokalt lagres det i lokal queue*/

var Local_queue = [4]int{0, 0, 1, 1}
var motor_direction_var int = 0

/*implement initialize*/
func init_elevator() {
	if elevio.GetFloor() == -1 {
		elevio.SetMotorDirection(-1)
	} else {
		elevio.SetMotorDirection(0)
	}
}

/*this section implements the functions that getthe motor started from idle*/
//helper functions to initialize the motor
func local_queue_check_above(sensor int) int {
	for floor := 0; floor < 4; floor++ {
		if 0 < Local_queue[floor] {
			if sensor < floor {
				return 1
			}
		}
	}
	return 0
}

func local_queue_check_below(sensor int) int {
	for floor := 0; floor < 4; floor++ {
		if 0 < Local_queue[floor] {
			if sensor > floor {
				return 1
			}
		}
	}
	return 0
}

func start_motor_from_idle() {

	if local_queue_check_above(elevio.GetFloor()) == 1 {
		set_motor_direction_variable(1)
		elevio.SetMotorDirection(1)
	}
	if local_queue_check_below(elevio.GetFloor()) == 1 {
		set_motor_direction_variable(-1)
		elevio.SetMotorDirection(-1)
	}
}

func check_order_at_floor() int {
	if Local_queue[elevio.GetFloor()] > 0 {
		return 1
	}
	return 0
}

func local_queue_check_for_saved_order() int {
	var i int = 0
	for i := 0; i < 4; i++ {
		if Local_queue[i] != 0 {
			i += Local_queue[i]
		}
	}
	if i > 0 {
		return 1
	} else {
		return 0
	}
}

func set_motor_direction_variable(direction int) {
	motor_direction_var = direction
	//possbily set a go routinefor
}
func get_motor_direction_variable() int {
	return motor_direction_var
}

func cab_button_at_floor() int {
	for floor := 0; floor < 4; floor++ {
		if 3 == Local_queue[floor] {
			if elevio.GetFloor() == floor {
				return 1
			}
		}
	}
	return 0
}

func down_button_at_floor() int {
	for floor := 0; floor < 4; floor++ {
		if 2 == Local_queue[floor] {
			if elevio.GetFloor() == floor {
				return 1
			}
		}
	}
	return 0
}

func up_button_at_floor() int {
	for floor := 0; floor < 4; floor++ {
		if 1 == Local_queue[floor] {
			if elevio.GetFloor() == floor {
				return 1
			}
		}
	}
	return 0
}

/*RUNNING --> DOOR*/
func check_if_correct_floor() int {
	if elevio.GetFloor() != -1 {
		if cab_button_at_floor() == 1 {
			return 1
		}
		if get_motor_direction_variable() == -1 {
			if down_button_at_floor() == 1 {
				return 1
			}
			if up_button_at_floor() == 1 {
				if local_queue_check_below(elevio.GetFloor()) == 0 {
					return 1
				}
			}
		}
		if get_motor_direction_variable() == 1 {
			if up_button_at_floor() == 1 {
				return 1
			}
			if down_button_at_floor() == 1 {
				if local_queue_check_above(elevio.GetFloor()) == 0 {
					return 1
				}
			}
		}
	}
	return 0
}

/*RUNNING -->MOTORFAILURE
implment the check enginer timer above threshold
*/

/*The logic part for what the door state should do*/

func open_door() {
	elevio.SetDoorOpenLamp(true)
}
func close_door() {
	elevio.SetDoorOpenLamp(false)
}

func local_queue_erase_floor_buttons() {
	if Local_queue[elevio.GetFloor()] != 0 {
		Local_queue[elevio.GetFloor()] = 0
		//posssibly send a go routine that updates the order handler module
	}
}

/*DOOR --> RUNNING*/
//start engine timer
//check that check_order comes not empty
func start_motor_from_door() {
	if local_queue_check_above(elevio.GetFloor()) == 1 && local_queue_check_below(elevio.GetFloor()) == 0 {
		set_motor_direction_variable(1)
	}
	if local_queue_check_above(elevio.GetFloor()) == 0 && local_queue_check_below(elevio.GetFloor()) == 1 {
		set_motor_direction_variable(-1)
	}
	if local_queue_check_for_saved_order() == 0 {
		set_motor_direction_variable(0)
	}
}

func local_queue_erase_all_buttons() {
	for floor := 0; floor < 4; floor++ {
		Local_queue[floor] = 0
	}
	//possily send a go routine that updates that the order handler module
}

/*functions that the cost function will need*/
func check_if_different_order_is_already_saved_at_floor(floor int, button int) int {
	if Local_queue[floor] == 0 {
		return 0
	} else {
		if Local_queue[floor] != button {
			return 1
		}
	}
	return 0
}

func FSM() {
	var STATE string = "INIT"
	Door_timer := time.Now()
	//engine_timer := timer.GetTime()
	for true {
		fmt.Println(STATE)
		switch STATE {
		case "INIT":
			if elevio.GetFloor() == -1 {
				elevio.SetMotorDirection(-1)
			} else {
				elevio.SetMotorDirection(0)
				STATE = "IDLE"
			}
			break
		case "IDLE":
			if check_order_at_floor() == 1 {
				Door_timer = time.Now()
				STATE = "DOOR"

				//set door timer
			}
			if check_order_at_floor() == 0 {
				start_motor_from_idle()
			}
			if get_motor_direction_variable() != 0 {
				//set enginer timer
				fmt.Println("hello world")
				STATE = "RUNNING"
			}

			break
		case "RUNNING":
			if check_if_correct_floor() == 1 {
				//door start timer
				Door_timer = time.Now()
				local_queue_erase_floor_buttons()
				elevio.SetMotorDirection(0)
				STATE = "DOOR"
			}
			// if motor failure
			break
		case "DOOR":

			open_door()
			if check_order_at_floor() == 1 {
				Door_timer = time.Now()
			}

			local_queue_erase_floor_buttons()

			if time.Now().Sub(Door_timer).Seconds() > 3 {
				close_door()
				start_motor_from_door()

				if get_motor_direction_variable() != 0 {
					value := get_motor_direction_variable()
					fmt.Println(value)

					switch value {
					case 1:
						elevio.SetMotorDirection(elevio.MD_Up)
						break
					case -1:
						elevio.SetMotorDirection(elevio.MD_Down)
						break
					}

					STATE = "RUNNING"
				} else {
					STATE = "IDLE"
				}
			}
			break
		case "MOTORFAILURE":
		}
	}
}

//if the cost function desides that the order shouldbe taken locally, it sends in the floor and button to this function so that local queue can be updated
func save_order_into_local_queue(floor int, button int) {
	if check_if_different_order_is_already_saved_at_floor(floor, button) == 1 {
		Local_queue[floor] = 3
	} else {
		Local_queue[floor] = button
	}
}

//need a function that constantly stores changes checks changes in the local_queue and updates the order matrix in order_handler

//The new state machine
type StateChannels struct {
	OrderComplete  chan int
	ArrivedAtFloor chan int
	//StateError 	chan error
	NewOrder chan config.Keypress
	Elevator chan config.Elev
}

func RunElevator(channel StateChannels) {
	elevator := config.Elev{
		State: config.Idle,
		Dir:   elevio.MD_Stop,
		Floor: elevio.GetFloor(),
		//Queue: [config.NumFloor][config.NumButtons]bool{},
	}

	DoorTimer := time.NewTimer(3 * time.Second)
	EngineFailureTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()
	EngineFailureTimer.Stop()

	orderCleared := false

	channel.Elevator <- elevator

	for {
		select {
		case newOrder := <-channel.NewOrder:

			if newOrder.Completed {
				elevator.Queue[newOrder.Floor][elevio.BT_HallUp] = false
				elevator.Queue[newOrder.Floor][elevio.BT_HallDown] = false
				//elevator.Queue[newOrder.Floor][elevio.BT_Cab] = false
				orderCleared = true

			} else {
				elevator.Queue[newOrder.Floor][newOrder.Btn] = true
			}

			switch elevator.State {
			case config.Idle:
				elevator.Dir = chooseDirection(elevator)
				if elevator.Dir == elevio.MD_Stop {
					elevator.State = config.DoorOpen
					elevio.SetDoorOpenLamp(true)
					DoorTimer.Reset(3 * time.Second)
					go func() { channel.OrderComplete <- newOrder.Floor }()
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{}

				} else {
					elevator.State = config.Moving
					elevio.SetMotorDirection(elevator.Dir)
					EngineFailureTimer.Reset(3 * time.Second)
				}
			case config.Moving:
			case config.DoorOpen:
				if elevator.Floor == newOrder.Floor {
					DoorTimer.Reset(3 * time.Second)
					go func() { channel.OrderComplete <- newOrder.Floor }()
					elevator.Queue[elevator.Floor] = [config.NumButtons]bool{}
				}
			case config.Undefined:
				fmt.Println("fatal error")
			}

		case elevator.Floor = <-channel.ArrivedAtFloor:
			fmt.Println("Arrived at floor")

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
			channel.Elevator <- elevator
		}
	}

}
