package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
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

func (p *RabbitPublisher) Publish(bindingKey string, json string, timestamp time.Time) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    timestamp,
		ContentType:  "application/json",
		Body:         []byte(json),
	}

	err := p.channel.Publish(RABBIT_EXCHANGE, bindingKey, false, false, msg)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return err
	}
	return nil
}
