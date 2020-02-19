package main

import "./network/bcast"
// import "./network/localip"
// import "./network/peers"
// import "flag"
import "fmt"
// import "os"
import "time"

func main(){
    
    // Make channels for listen and send
    numRx := make(chan int)
    numTx := make(chan int)
    
    // First listen, then send
    go bcast.Receiver(16569, numRx)
    
    var num int = 0 // How far we have counted
    var lastNum int = 0 // Previous number
    var t time.Time = time.Now() // Current time
    fmt.Println("Listening for counter...")
    
    listen:
    for {
        select{
        case num = <-numRx:
            t = time.Now() //update for each recieve
            if(num != lastNum){
                fmt.Println("External counter:", num)
            }
        default:
            if(time.Since(t) > 2*time.Second){
                fmt.Println("External counting stopped. Resuming count locally.")
                break listen // stop listening
            }
        }
    }
    
    // CREATE BACKUP PROGRAM
    
    
    for{
        go bcast.Transmitter(16569, numTx)
        fmt.Println("Local counter:", num)
        num++
        numTx <- num
        time.Sleep(500*time.Millisecond)
    }
}