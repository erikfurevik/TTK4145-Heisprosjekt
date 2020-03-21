package main

import (
	"./config"
	"./elevio"
	"./fsm"
	bcast "./network/bcast"
	"fmt"
	"time"
)

func main() {
	elevio.Init("localhost:15657", config.NumFloor)
	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	//receiveOrders := make(chan elevio.ButtonEvent)
	//receiveFloors := make(chan int)
	channels := fsm.StateChannels{
		OrderComplete:  make(chan int),
		Elevator:       make(chan config.Elev),
		NewOrder:       make(chan elevio.ButtonEvent),
		ArrivedAtFloor: make(chan int),
	}


	go elevio.PollButtons(channels.NewOrder)
	go elevio.PollFloorSensor(channels.ArrivedAtFloor)

	go fsm.RunElevator(channels)

	//Under her finner man kode og kommentarer for å teste og forstå transmitter og reciever litt mer

	// [heisId | et1 | et2 | et3 | et4] , bare en veldig enkel array for å se at det kunne sendes
	var orders =[5] int{1,0,0,1,0}
	var recievedMsg [5]int
	id := 1
	port :=12000 //Har bare valgt en random port for å teste kode.
	messageTx := make(chan [5]int)
	messageRx := make(chan [5]int)

	// Transmitter og receiver kan ikke sende structs som de er nå, men det er mulig å endre ved annerledes bruk av "reflect"
	// Hva skal egentlig sendes? For å sende arrays går ihvertfall fint. 
	go bcast.Transmitter(port, messageTx) //port, melding som sendes.
	go bcast.Receiver(port,messageRx) //port, motatt melding
	

	for {
		messageTx <- orders
		select{
		case recievedMsg = <- messageRx:
			if recievedMsg[0] != id{  //testet med å kjøre koden på flere terminaler med ulik id
				fmt.Println("Heis ",recievedMsg[0], " har bestillinger i: ", recievedMsg[1:5])
			}
		}
		time.Sleep(500*time.Millisecond)
		//fmt.Println(timer.CheckTime(t))
	}

}
