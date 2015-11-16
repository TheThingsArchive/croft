package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/thethingsnetwork/server-shared"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type MQTTPublisher struct {
	cli *client.Client
}

func (mq *MQTTPublisher) Configure() error {
	//NOOP
	return nil
}

func (mq *MQTTPublisher) Publish(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %s", err.Error())
		return err
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

	// Publish a message.
	return mq.cli.Publish(&client.PublishOptions{
		// QoS is the QoS of the PUBLISH Packet.
		QoS: mqtt.QoS0,
		// Retain is the Retain of the PUBLISH Packet.
		Retain: true,
		// TopicName is the Topic Name of the PUBLISH Packet.
		TopicName: []byte(routingKey),
		// Message is the Application Message of the PUBLISH Packet.
		Message: []byte(body),
	})
}

func ConnectMQTTPublisher() (Publisher, error) {
	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// // Terminate the Client.
	// defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  fmt.Sprintf("%s:1883", os.Getenv("MOSQUITTO_URI")),
		ClientID: []byte("croft"),
	})

	if err != nil {
		return &MQTTPublisher{nil}, err
	}
	log.Printf("Connected to %s", fmt.Sprintf("%s:1883", os.Getenv("MOSQUITTO_URI")))

	return &MQTTPublisher{cli}, nil
}
