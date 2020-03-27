package ElevatorController

import (
    "../config"
	"../elevio"
)

//Likewise if the 
func costFunction(Local_ID int, LocalOrder elevio.ButtonEvent, elevatorList [config.NumElevator]config.Elev, elevatorOnline [config.NumElevator]bool) int {
    if LocalOrder.Button == elevio.BT_Cab {
        return Local_ID
    }
    var CostArray[config.NumElevator]int

    for elev := 0; elev < config.NumElevator; elev++ {
        cost := LocalOrder.Floor - elevatorList[elev].Floor

        if cost == 0 && elevatorList[elev].State != config.Moving {
            return elev
        }
        
        if cost < 0 {
            cost = -cost
            if elevatorList[elev].Dir == elevio.MD_Up {
                cost += 3
            }
        } else if cost > 0 {
            if elevatorList[elev].Dir == elevio.MD_Down {
            cost += 3
            }
        }

        if cost == 0 && elevatorList[elev].State == config.Moving {
            cost += 4
        }

        if elevatorList[elev].State == config.DoorOpen {
            cost ++
        
        }
        CostArray[elev] = cost
    }

    maxCost := 1000;
    var bestElev int
    for elev := 0; elev < config.NumElevator; elev++{
        if CostArray[elev] < maxCost && elevatorOnline[elev] && elevatorList[elev].State != config.Undefined{
            bestElev = elev
            maxCost = CostArray[elev]
        }
    } 

    return bestElev
    
}
