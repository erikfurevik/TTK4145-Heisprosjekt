package main

import "fmt"
import "time"

func main(){
    
    NetTx := make(chan ElevatorData)    // Transmit to network
    NetRx := make(chan ElevatorData)    // Receive from network
    
    go network(NetTx, NetRx)            // Start network module
    
    info := ElevatorData{"10.22.72.183", 0} // create elevator struct
    go func(){
        var p int = 0
        for{
            info.Timestamp++
            p++
            NetTx<- info    // Send to network
            time.Sleep(time.Second)
        }
    }()
    
    for{
        select{
        case r:= <-NetRx:   // Message from network
            fmt.Println(r)
        }
    }
}
