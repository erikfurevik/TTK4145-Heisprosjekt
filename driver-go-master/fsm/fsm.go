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


//helper functions to initialize the motor 
func local_queue__check_above(sensor int)int{
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

func erase_all_buttons(){
	for  var floor int = 0; floor < 4; floor++{
		local_queue{floor] = 0}
	}
}

func local_queue_erase_floor_buttons(floor int){
	local_queue[floor] = 0
} 


//The logic part whether the elavator should stop at floor


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