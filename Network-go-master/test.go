package main

import "fmt"
import "time"

func main(){
    
    NetTx := make(chan ElevatorData)    // Transmit to network
    NetRx := make(chan ElevatorData)    // Receive from network
    
    go network(NetTx, NetRx)            // Start network module
    
    info := ElevatorData{"Hi", 0}       // create elevator struct
    
    for{
        info.Timestamp++
        NetTx<- info                    // Send to network
        time.Sleep(time.Second*1)
        select{
        case r:= <-NetRx:               // Message from network
            fmt.Println(r)
        }
    }
    
}
