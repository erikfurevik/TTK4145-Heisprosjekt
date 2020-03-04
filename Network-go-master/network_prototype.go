package main
import (
	"./network/bcast"
	"./network/localip"
	// "./network/peers"
	// "flag"
	"fmt"
	// "os"
	// "strconv"
)

var NUMBER_OF_FLOORS int = 4

func network(NetTx <-chan ElevatorData, NetRx chan<- ElevatorData){
	
	// Get local IP adress
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error")
	}
	// We could convert IP like this:
	// splitIP := strings.Split(localIP, ",")
	// intIP, err := strconv.Atoi(splitIP[3])
	
	dataRx := make(chan ElevatorData)
	dataTx := make(chan ElevatorData)
	
	
	go bcast.Transmitter(16731, dataTx)	// start broadcaster
	go bcast.Receiver(16731, dataRx)	// start broadcast receiver
	
	// This sends data to network. Test only.
	// Should replace with real data from main.

	
	for{
		select{
		case r := <-dataRx:
			if r.ID == localIP{
				// Sent from local computer
				fmt.Println("Ignore")
			}
			NetRx <- r
			
		case s := <-NetTx:
			dataTx<- s
		}
	}
}
