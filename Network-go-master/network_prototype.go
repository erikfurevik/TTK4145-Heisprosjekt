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
	"reflect"
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
var elevData map[string]ElevatorData	// map of all data
var rxCounter map[string]int			// map of times received data

func network(NetTx <-chan ElevatorData, NetRx chan<- ElevatorData){
	
	// Get local IP adress
	localAdress, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error, no connection!")
		return
	}
	localID := localAdress
	fmt.Println("Local ID:", localID)
	// localID := strings.Split(localAdress, ".")[3]
	// Can be converted to int like this:
	// intIP, err := strconv.Atoi(splitIP[3])
	
	peers = make(map[string]time.Time)
	elevData = make(map[string]ElevatorData)
	rxCounter = make(map[string]int)
	
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
		
			// Checks for correct datatype
			if(reflect.TypeOf(r).String() != "main.ElevatorData"){
				fmt.Println("Wrong datatype received")
				break
			}
			
			// Handles message errors
			if elevData[r.ID] != r {
				rxCounter[r.ID] = 0
				fmt.Println("New message")
			}
			elevData[r.ID] = r
			rxCounter[r.ID]++
			if rxCounter[r.ID] < 4 {
				fmt.Println("Waiting for confirmation")
				break	// Message may be faulty
			}
			
			if peers[r.ID].IsZero() {
				if r.ID == localID {
					fmt.Println("We're back online:")
					// regained connection to network
				} else {
					fmt.Println("New peer:")
					NetRx <- r
				}
			} else if r.ID == localID{
				fmt.Println("Message from self:")
				// ignore
			} else {
				fmt.Println("Update from peer:")
				NetRx <- r
			}
			peers[r.ID] = time.Now() // Update receive time
			
		case s := <-NetTx:
			// sends to network
			dataTx <- s
		}
	}
}
