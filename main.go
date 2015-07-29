package main

import (
	"log"
)

var (
	publisher Publisher
)

func main() {
	log.Print("Croft is ALIVE")

	connectPublisher()
	startUDPServer(1700)
}

func connectPublisher() {
	var err error
	publisher, err = ConnectRabbitPublisher()
	if err != nil {
		log.Fatal(err)
	}

	err = publisher.Configure()
	if err != nil {
		log.Fatal(err)
	}
}
