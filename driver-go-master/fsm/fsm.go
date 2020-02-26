package fsm 


import (
	"fmt"
	"../elevio"

)

const floors int

/*value 1 = up
value 2 = down
value 3 = cab
*/

/*så fort kost funksjonen bergner om en ordre skal som skal bli tatt lokalt lagres det i lokal queue*/

var local_queue[floors] int
var motor_direction_var int


/*this section implements the functions that getthe motor started from idle*/
//helper functions to initialize the motor 
func local_queue_check_above(sensor int)int{
	for var floor int = 0; floor < 4; floor++ {
		if 0 < local_queue[floor]{
			if sensor < floor{
				return 1
			}
		}
	}
	return 0
}

func local_queue_check_below(sensor int) int{
	for var floor int = 0; floor < 4; floor++{
		if 0 < local_queue[floor] {
			if sesnor > floor {
				return 1
			}
		}
	}
	return 0
}

func local_queue_check_for_saved_order()int{
	var i int = 0
	for int i = 0 i < 4; i++ {
		if local_queue[i] != 0{
			i += local_queue[i]	
		}
	}
	if i > 0 {
		return 1
	}
	else {
		return 0
	}
}


func set_motor_direction_variable(direction int){
	motor_direction_var = direction
	//possbily set a go routine
}

func get_motor_direction_variable()int {
	return motor_direction_var
}

/*IDLE -->RUNNING*/
func start_motor_from_idle(){
	if local_queue_check_for_saved_order(){
		if local_queue_check_above(elevio.getFloor()){
			set_motor_direction_variable(1)
			elevio.setMotorDirection(1)
		}
		if local_queue_check_below(elevio.getFloor()){
			set_motor_direction_variable(-1)
			elevio.setMotorDirection(-1)
		}
	}
}

/*IDLE -->DOOR*/
func check_order_at_floor(){
	for var floor int = 0; floor < 4; floor++{
		if local_queue[floor] > 0{
			if floor == elevio.getFloor(){
				return 1
			}
		}
	} 
}


/*IDLE -->RUNNING 
implemnt the engine timer start 
*/

/*IDLE -->DOOR
implement the door start timer
*/



/*The logic part whether the elavator should stop at floor*/


//help functions
func up_button_at_floor()int{
	for var floor int = 0; floor < 4; floor++{
			if local_queue[floor] == 1{
				if floor == elevio.getFloor(){
					return 1
				} 
			}
	} 
	return 0
}


func down_button_at_floor()int{
	for var floor int = 0; floor < 4; floor++{
		if local_queue[floor] = 2{
			if local_queue[floor] = elevio.getFloor(){
				return 1
			}
		}
	}
	return 0
}

func cab_button_at_floor()int{
	for var floor int = 0; floor < 4; floor++{
		if local_queue[floor] = 3{
			if local_queue[floor] = elevio.getFloor(){
				return 1
			}
		}
	}
	return 0
}



/*RUNNING --> DOOR*/
func check_if_correct_floor(){
	if elevio.getFloor() != -1{	
		if cab_button_at_floor(){
			return 1
		}
		if get_motor_direction_variable() == -1{
			if down_button_at_floor() {
				return 1
			}
			if up_button_at_floor(){
				if !local_queue_check_below(elevio.getFloor()){
					return 1
				}
			}
		}
		if get_motor_direction_variable() = 1{
			if up_button_at_floor(){
				return 1
			}
			if down_button_at_floor(){
				if !local_queue_check_above(elevio.getFloor()){
					return 1
				}
			}
		}
	}
}


/*RUNNING -->DOOR 
implment the door start timer*/


/*RUNNING -->MOTORFAILURE
implment the check enginer timer above threshold
*/


/*The logic part for what the door state should do*/

func open_door(){
	elevio.SetDoorOpenLamp(true)
}
func close_door(){
	elevio.SetDoorOpenLamp(false)
}

func local_queue_erase_floor_buttons(){
	local_queue[elevio.getFloor()] = 0
	//posssibly send a go routine that updates the order handler module
} 

/*DOOR -->IDLE ||RUNNING
implement the check door timer above threshold*/


/*DOOR -->IDLE*/
//check for check_order comes empty
//set_motor_direction_variable(0)

/*DOOR --> RUNNING*/
//start engine timer
//check that check_order comes not empty
func start_motor_from_door(){
	if local_queue_check_above(elevio.getFloor()) && !local_queue_check_below(elevio.getFloor()){
		set_motor_direction_variable(1)
	}
	if !local_queue_check_above(elevio.getFloor()) && local_queue_check_below(elevio.getFloor()){
		set_motor_direction_variable(-1)
	}
}





func erase_all_buttons(){
	for  var floor int = 0; floor < 4; floor++{
		local_queue{floor] = 0}
	}
	//possily send a go routine that updates that the order handler module
}



/*functions that the cost function will need*/
func check_if_different_order_is_already_saved_at_floor(floor int, button int)int{
	if local_queue{floor] == 0 {
		return 0
	}
	else {
		if local_queue[floor] != button {
			return 1
		}
	}
}


//if the cost function desides that the order shouldbe taken locally, it sends in the floor and button to this function so that local queue can be updated
func save_order_into_local_queue(floor int, button int){
	if check_if_different_order_is_already_saved_at_floor(floor, button){
		local_queue[floor] = 3;
	}
	else{
		local_queue[floor] = button
	}
}


//need a function that constantly stores changes checks changes in the local_queue and updates the order matrix in order_handler