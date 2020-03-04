package main
import (
	"./network/bcast"
	"./network/localip"
	// "./network/peers"
	// "flag"
	"fmt"
	"time"
	// "strings"
	// "os"
	// "strconv"
)


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

var NUMBER_OF_FLOORS int = 4			// number of floors
var NUMBER_OF_ELEVATORS int = 10 		// number of elevators (max)
var peers map[string]time.Time	// map of all peers



func network(NetTx <-chan ElevatorData, NetRx chan<- ElevatorData){
	
	// Get local IP adress
	localAdress, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error, no connection!")
		return
	}
	localID := localAdress
	// localID := strings.Split(localAdress, ".")[3]
	// Can be converted to int like this:
	// intIP, err := strconv.Atoi(splitIP[3])
	fmt.Println(localID)
	
	peers = make(map[string]time.Time)
	peers[localID] = time.Now()
	
	dataRx := make(chan ElevatorData)
	dataTx := make(chan ElevatorData)
	
	go bcast.Transmitter(16731, dataTx)	// start broadcaster
	go bcast.Receiver(16731, dataRx)	// start broadcast receiver
	
	//TODO:
	// Keep track of all online elevators
	// One timer for each elevator
	// When timer runs out, remove elevator
	// If new elevator found, add elevator
	
	
	
	for{
		select{
		case r := <-dataRx:	// Message received
			peers[r.ID] = time.Now() // Update time
			if r.ID == localID {
				
			} else if  {
				
			}
			NetRx <- r
	
		case s := <-NetTx:
			dataTx<- s
		}
	}
}
