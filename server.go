package main

import (
	"fmt"
	"github.com/thethingsnetwork/croft/lora"
	"github.com/thethingsnetwork/server-shared"
	"log"
	"net"
)

func readUDPMessages(port int, messages chan interface{}) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to resolve address: %s", err.Error())
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %s", err.Error())
	}
	defer listener.Close()

	conn := lora.NewConn(listener)
	for {
		msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err.Error())
			continue
		}
		log.Printf("Parsed message: %#v", msg.Header)

		go handleMessage(msg, messages)
	}
}

func handleMessage(msg *lora.Message, messages chan interface{}) {
	switch msg.Header.Identifier {
	case lora.PUSH_DATA:
		publishPushMessagePayloads(msg.Payload.(lora.PushMessagePayload), messages)
	}

	err := msg.Ack()
	if err != nil {
		log.Printf("Failed to send ACK: %s", err.Error())
	}
}

func publishPushMessagePayloads(payload lora.PushMessagePayload, messages chan interface{}) {
	if payload.Stat != nil {
		stat, err := convertStat(payload.Stat)
		if err != nil {
			log.Printf("Failed to convert Stat: %s", err.Error())
		}
		messages <- stat
	}

	if payload.RXPK != nil {
		for _, rxpk := range payload.RXPK {
			packet, err := convertRXPK(rxpk)
			if err != nil {
				log.Printf("Failed to convert RXPK: %s", err.Error())
				continue
			}
			messages <- packet
		}
	}
}

func convertStat(stat *lora.Stat) (*shared.GatewayStatus, error) {
	return &shared.GatewayStatus{
		Time:              stat.Time,
		Latitude:          &stat.Lati,
		Longitude:         &stat.Long,
		Altitude:          &stat.Alti,
		RxCount:           &stat.Rxnb,
		RxOk:              &stat.Rxok,
		RxForwarded:       &stat.Rxfw,
		AckRatio:          &stat.Ackr,
		DatagramsReceived: &stat.Dwnb,
		DatagramsSent:     &stat.Txnb,
	}, nil
}

func convertRXPK(rxpk *lora.RXPK) (*shared.RxPacket, error) {
	return &shared.RxPacket{
		Time: rxpk.Time,
		Data: rxpk.Data,
	}, nil
}
