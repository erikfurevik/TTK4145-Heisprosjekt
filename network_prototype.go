package main
import (
	"./Network-go-master/network/bcast"
	// "./network/localip"
	// "./network/peers"
	// "flag"
	"fmt"
	// "os"
	"time"
)

var NUMBER_OF_FLOORS int = 4

type ElevatorData struct {
    // ID          int
    TimeStamp   int
	// State		int
	// Location	int
	// Direction	int
	// RequestsUp [NUMBER_OF_FLOORS-1]bool
	// RequestsDown [NUMBER_OF_FLOORS-1]bool
	// RequestsCab [NUMBER_OF_FLOORS]bool
}

func main(){

	dataRx := make(chan ElevatorData)
	dataTx := make(chan ElevatorData)

	sendData := make(chan ElevatorData)

	go bcast.Transmitter(16570, dataTx)
	go bcast.Receiver(16570, dataRx)

	// This sends data to network. Test only.
	// Should replace with real data from main.
	go func(){
		info := ElevatorData{0}
		for{
			info.TimeStamp++
			sendData<- info
			time.Sleep(time.Second*1)
		}
	}()

	for{
		select{
		case r := <-dataRx:
			fmt.Println(r)
		case s := <-sendData:
			dataTx<- s
		}
	}
}
