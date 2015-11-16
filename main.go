package main

import (
	"log"
)

var (
	publishers []Publisher
)

func main() {
	log.Print("Croft is ALIVE")

	rabbitPublisher, err := connectRabbitPublisher()
	if err != nil {
		log.Fatalf("Failed to connect rabbit publisher: %s", err.Error())
	}

	mqttPublisher, err := connectMQTTPublisher()
	if err != nil {
		log.Fatalf("Failed to connect MQTT publisher: %s", err.Error())
	}

	messages := make(chan interface{})
	go readUDPMessages(1700, messages)
	for msg := range messages {
		err = rabbitPublisher.Publish(msg)
		if err != nil {
			log.Printf("Failed to publish message to rabbit %#v: %s", msg, err.Error())
		}
		err = mqttPublisher.Publish(msg)
		if err != nil {
			log.Printf("Failed to publish message to mqtt %#v: %s", msg, err.Error())
		}
	}
}

func connectRabbitPublisher() (Publisher, error) {
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

func connectMQTTPublisher() (Publisher, error) {
	publisher, err := ConnectMQTTPublisher()
	if err != nil {
		return nil, err
	}

	err = publisher.Configure()
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
