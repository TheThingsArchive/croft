package main

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	rabbitConn    *amqp.Connection
	rabbitChannel *amqp.Channel
)

func main() {
	log.Print("Croft is ALIVE")

	var err error
	rabbitConn, err = EnsureRabbitConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()
	rabbitChannel, err = ConfigureRabbit()
	if err != nil {
		log.Fatal(err)
	}
	err = PublishRabbitMessage()
	if err != nil {
		log.Fatal(err)
	}
	go StartUDPServer(1700)
	ServeHTTPOverview()
}
