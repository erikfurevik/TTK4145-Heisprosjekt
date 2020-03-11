package Order_handler

func costFunction(request Keypress, elevatorList [NumElevator]Elev, id int, elevatorOnline [NumElevator]bool) int {
    if order.Btn == elevio.BT_Cab {
        return id
    }

    minCost := numButtons*NumFloor*NumElevator

    for elev := 0; elev < NumElevator; elevator++ {
        if !elevatorOnline[elev] {
            continue
        }

        cost := request.Floor - elevatorList[elev].Floor

        if cost == 0 && elevatorList[elev].State != Moving {
            return elev
        }

        if cost < 0 {
            cost = -cost
            if elevatorList[elev].Dir == MD_Up {
                cost += 3
            }
        } else if cost > 0 {
            if elevatorList[elev].Dir == MD_Down {
            cost += 3
            }
        }

        if cost == 0 && elevatorList[elev].State == Moving {
            cost += 4
        }

        if elevatorList[elev].State == DoorOpen {
            cost += 1
        }

        if cost < minCost {
            minCost = cost
            bestElev = elev
        }
    }
    return bestElev
}
