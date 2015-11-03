package main

import (
	"log"
)

func main() {
	log.Print("Croft is ALIVE")

	publisher, err := connectPublisher()
	if err != nil {
		log.Fatalf("Failed to connect publisher: %s", err.Error())
	}

	messages := make(chan interface{})
	go readUDPMessages(1700, messages)
	for msg := range messages {
		err = publisher.Publish(msg)
		if err != nil {
			log.Printf("Failed to publish message %#v: %s", msg, err.Error())
		}
	}
}

func connectPublisher() (Publisher, error) {
	publisher, err := ConnectRabbitPublisher()
	if err != nil {
		return nil, err
	}

	err = publisher.Configure()
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
