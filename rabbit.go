package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/thethingsnetwork/server-shared"
)

const (
	RABBIT_ATTEMPTS = 20
	RABBIT_EXCHANGE = "messages"
)

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func ConnectRabbitPublisher() (Publisher, error) {
	var err error
	for i := 0; i < RABBIT_ATTEMPTS; i++ {
		uri := os.Getenv("AMQP_URI")
		conn, err := amqp.Dial(uri)
		if err != nil {
			log.Printf("Failed to connect: %s", err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			publisher := &RabbitPublisher{conn, nil}
			log.Printf("Connected to %s", uri)
			return publisher, nil
		}
	}
	return nil, err
}

func (p *RabbitPublisher) Configure() error {
	c, err := p.conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		return err
	}

	err = c.ExchangeDeclare(RABBIT_EXCHANGE, "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		return err
	}

	p.channel = c
	return nil
}

func (p *RabbitPublisher) Publish(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %s", err.Error())
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	}

	var routingKey string
	switch data.(type) {
	case *shared.GatewayStatus:
		routingKey = "gateway.status"
	case *shared.RxPacket:
		routingKey = "gateway.rx"
	default:
		return errors.New("Invalid type to publish")
	}

	err = p.channel.Publish(RABBIT_EXCHANGE, routingKey, false, false, msg)
	if err != nil {
		log.Printf("Failed to publish: %s", err.Error())
		return err
	}
	return nil
}
