package controller

import (
    "../config"
	"../elevio"
)

func costFunction(Local_ID int, LocalOrder elevio.ButtonEvent, elevatorList [config.NumElevator]config.Elev, elevatorOnline [config.NumElevator]bool) int {
    if LocalOrder.Button == elevio.BT_Cab {
        return Local_ID
    }

    minCost := config.NumButtons*config.NumFloor*config.NumElevator
    bestElev := Local_ID
    var secondBest int 

    for elev := 0; elev < config.NumElevator; elev++ {
        if !elevatorOnline[elev] {
            continue
        }
        secondBest = elev
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
        if cost < minCost {
            minCost = cost
            bestElev = elev
        }
    }
    if bestElev == Local_ID && elevatorOnline[Local_ID] == false {
        return secondBest
    }else{
        return bestElev
    }
}
