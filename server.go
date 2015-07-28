package main

import (
	"fmt"
	"github.com/thethingsnetwork/croft/lora"
	"log"
	"net"
	"time"
)

var lc *lora.Conn

func startUDPServer(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	CheckError(err)

	listener, err := net.ListenUDP("udp", addr)
	CheckError(err)
	defer listener.Close()

	lc = lora.NewConn(listener)
	for {
		msg, err := lc.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err.Error())
			continue
		}
		log.Printf("Parsed message: %#v", msg.Header)

		if msg.Header.Identifier == lora.PUSH_DATA {
			for _, field := range *msg.Payload {
				log.Printf("Publishing message %s: %s", field.Key, field.Value)
				err = publisher.Publish(field.Key, field.Value, time.Now())
				if err != nil {
					log.Printf("Error publishing %s: %s", field.Key, err.Error())
					continue
				}
				log.Printf("Published %s", field.Key)
			}
		}

		err = msg.Ack()
		if err != nil {
			log.Printf("Error sending ACK: %s", err.Error())
		}
	}
}
