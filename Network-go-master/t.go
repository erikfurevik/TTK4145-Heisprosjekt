package main

import "fmt"
// import "time"
import "reflect"

type ElevatorData struct {	// state of one elevator
    ID          string
    Timestamp   int
	// State       int
	// Location	   int
	// Direction   int
	// RequestsUp [NUMBER_OF_FLOORS-1]bool
	// RequestsDown [NUMBER_OF_FLOORS-1]bool
	// RequestsCab [NUMBER_OF_FLOORS]bool
}

func main(){
    info := ElevatorData{"10.22.72.183", 0} // create elevator struct
    p := reflect.TypeOf(info).String()
}
