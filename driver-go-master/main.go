package main

import (
	"fmt"
	"strconv"
	"time"

	"./config"
	"./elevio"
	"./fsm"
	bcast "./network/bcast"
	peers "./network/peers"
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
	// Burde selvsagt ikke hardkode antall elementer i tilfelle, men gadd ikke tenkte på det nå siden vi må
	// endre på hva som skal sendes uansett
	id := 1 //Må velge en måte å tildele id på
	var (
		orders      = [5]int{id, 0, 0, 1, 0}
		recievedMsg [5]int
		peerChange  peers.PeerUpdate
	)

	id_string := strconv.Itoa(id)
	port := 12000 //Har bare valgt en random port for å teste kode.

	messageTx := make(chan [5]int)            // Kanal for å sende melding
	messageRx := make(chan [5]int)            // kanal for å motta melding
	enableTx := make(chan bool)               // Kanal for om vi kan transmitte id eller ikke
	peerUpdate := make(chan peers.PeerUpdate) // brukes for å oppdatere hvor mange peers vi har, men litt usikker på hvordan

	// Transmitter og receiver kan ikke sende structs som de er nå, men det er mulig å endre ved annerledes bruk av "reflect"
	// Hva skal egentlig sendes? For å sende arrays går ihvertfall fint.
	go bcast.Transmitter(port, messageTx)           // Broadcaster melding
	go bcast.Receiver(port, messageRx)              // Mottar melding
	go peers.Transmitter(port, id_string, enableTx) //Transmitter id for å si man er på nettet
	go peers.Receiver(port, peerUpdate)             // funksjonen sjekker hvem som er på nettet og oppdaterer oversikten.

	for {
		messageTx <- orders
		select {
		case recievedMsg = <-messageRx:
			if recievedMsg[0] != id { //testet med å kjøre koden på flere terminaler med ulik id
				fmt.Println("Heis ", recievedMsg[0], " har bestillinger i: ", recievedMsg[1:5])
			}
		case peerChange = <-peerUpdate: //denne får in både id + orders tydeligvis
			fmt.Println(peerChange)
		}
		time.Sleep(500 * time.Millisecond)
		//fmt.Println(timer.CheckTime(t))
	}

}
