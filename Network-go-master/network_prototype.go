package main
import (
	"./network/bcast"
	"./network/localip"
	// "./network/peers"
	// "flag"
	"fmt"
	// "os"
	"time"
	// "strconv"
)

var NUMBER_OF_FLOORS int = 4

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
	
	sendData := make(chan ElevatorData)
	
	go bcast.Transmitter(16731, dataTx)
	go bcast.Receiver(16731, dataRx)
	
	// This sends data to network. Test only.
	// Should replace with real data from main.
	go func(){
		info := ElevatorData{localIP, 0}
		for{
			info.TimeStamp++
			sendData<- info
			time.Sleep(time.Second*1)
		}
	}()
	
	for{
		select{
		case r := <-dataRx:
			if r.ID == localIP{
				// Sent from local computer
				fmt.Println("Ignore")
			}else{
				// From different computer
				fmt.Println(r)
			}
		case s := <-sendData:
			dataTx<- s
		}
	}
}
