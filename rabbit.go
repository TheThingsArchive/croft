package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func EnsureRabbitConnection() (*amqp.Connection, error) {
	// Connects opens an AMQP connection from the credentials in the URL.
	var err error
	for i := 0; i < 20; i++ {
		conn, err := amqp.Dial(os.Getenv("AMQP_URI"))
		if err != nil {
			log.Print("Couldn't get rabbit connection")
			log.Print(err.Error())
			time.Sleep(time.Duration(2) * time.Second)
			log.Print("Retrying.....")
		} else {
			return conn, nil
		}
	}
	log.Print("Got rabbit connection")
	return nil, err
}

func ConfigureRabbit() (*amqp.Channel, error) {
	c, err := rabbitConn.Channel()
	if err != nil {
		log.Printf("channel.open: %s", err)
		return nil, err
	}

	// We declare our topology on both the publisher and consumer to ensure
	// they
	// are the same.  This is part of AMQP being a programmable messaging
	// model.
	//
	// See the Channel.Consume example for the complimentary declare.
	err = c.ExchangeDeclare("messages", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("exchange.declare: %v", err)
		return nil, err
	}
	return c, nil
}

func PublishRabbitMessage() error {
	// Prepare this message to be persistent.  Your publishing
	// requirements may
	// be different.
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte("Go Go AMQP!"),
	}
	// This is not a mandatory delivery, so it will be
	// dropped if there are no
	// queues bound to the logs exchange.
	err := rabbitChannel.Publish("messages", "stat", false, false, msg)
	if err != nil {
		// Since publish is asynchronous this can
		// happen if the network connection
		// is reset or if the server has run out
		// of resources.
		log.Printf("basic.publish: %v", err)
		return err
	}
	return nil
}
