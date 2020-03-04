package main

import "../network/bcast"
// import "./network/localip"
// import "./network/peers"
// import "flag"
import "fmt"
import "os/exec"
import "time"

func main(){
    
    // Channels for listen and send
    numRx := make(chan int)
    numTx := make(chan int)
    
    go bcast.Receiver(16569, numRx)
    
    var num int = 0 // How far we have counted
    var lastNum int = 0 // Previous number
    var t time.Time = time.Now() // Current time
    fmt.Println("Listening for counter...")
    
    listen:
    for { // Listening
        select{
        case num = <-numRx:
            t = time.Now() //update for each recieve
            if(num != lastNum){
                fmt.Println("External counter:", num)
            }
        default:
            if(time.Since(t) > 2*time.Second){ // Noone else is counting
                if(num != 0){
                    fmt.Println("External counting stopped,")
                }
                if(num == 0){
                    fmt.Println("No counting detected,")
                }
                break listen // Stop listening
            }
        }
    }
    
    fmt.Println("Creating clone...")
    go exec.Command("gnome-terminal", "-x", "go", "run", "exercise_6.go").Run()
    
    go bcast.Transmitter(16569, numTx)
    
    // Give time to setup clone and transmitter
    time.Sleep(1*time.Second)
    
    fmt.Println("Starting counter:")
    for{ // Counting
        num++
        fmt.Println("Local counter:", num)
        numTx <- num
        time.Sleep(500*time.Millisecond)
    }
}
