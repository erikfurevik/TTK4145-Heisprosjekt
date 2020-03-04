package main

import "fmt"
import "time"

type ElevatorData struct {
    ID          string
    TimeStamp   int
	// State		int
	// Location	int
	// Direction	int
	// RequestsUp [NUMBER_OF_FLOORS-1]bool
	// RequestsDown [NUMBER_OF_FLOORS-1]bool
	// RequestsCab [NUMBER_OF_FLOORS]bool
}

func main(){
    
    NetTx := make(chan ElevatorData)    // Ask network to send this
    NetRx := make(chan ElevatorData)    // Receive message from network
    
    go network(NetTx, NetRx)            // Start network module
    
    info := ElevatorData{"Hi", 0}       // create elevator struct
    
    for{
        info.TimeStamp++
        NetTx<- info                    // Send to network
        time.Sleep(time.Second*1)
        select{
        case r:= <-NetRx:               // Message from network
            fmt.Println(r)
        }
    }
    
}
