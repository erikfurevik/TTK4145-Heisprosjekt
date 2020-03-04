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

const NUMBER_OF_FLOORS int = 4			// number of floors
const NUMBER_OF_ELEVATORS int = 10 		// number of elevators (max)
const NETWORK_TIMEOUT time.Duration = time.Second*5 // network timeout
const NETWORK_POLLRATE time.Duration = 20 * time.Millisecond

var peers map[string]time.Time			// map of all peers

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
	
	peers = make(map[string]time.Time)
	
	dataRx := make(chan ElevatorData)
	dataTx := make(chan ElevatorData)
	
	go bcast.Transmitter(16731, dataTx)	// start broadcaster
	go bcast.Receiver(16731, dataRx)	// start broadcast receiver
	
	
	// NETWORK WATCHDOG:
	// deletes peers from map if time runs out
	// pollrate and timeout must be tuned
	go func(){
		for{
			time.Sleep(NETWORK_POLLRATE)
			for k, v := range peers {
				if time.Since(v) > NETWORK_TIMEOUT {
					delete(peers, k)
					fmt.Println("Deleted", k)
					break
				}
			}
		}
	}()
	
	// NETWORK SEND AND RECEIVE:
	// first part reads from network,
	// second part sends (passes on message)
	for{
		select{
		case r := <-dataRx:	// Message received
			if peers[r.ID].IsZero() {
				if r.ID == localID {
					fmt.Println("We're back online:")
					// regained connection to network
				} else {
					fmt.Println("New peer:")
					// transmit data
				}
			} else if r.ID == localID{
				fmt.Println("Message from self:")
				// ignore
			} else {
				fmt.Println("Update from peer:")
				// put data into storage
			}
			peers[r.ID] = time.Now() // Update receive time
			NetRx <- r // Pass on data
	
		case s := <-NetTx:
			dataTx<- s
		}
	}
}
