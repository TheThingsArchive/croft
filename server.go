package main

import (
	"fmt"
	"github.com/thethingsnetwork/croft/lora"
	"log"
	"net"
)

var lc *lora.Conn

func StartUDPServer(port int) {
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	CheckError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	lc = lora.NewConn(ServerConn)

	for {
		msg, err := lc.ReadMessage()
		if err != nil {
			continue
		}
		log.Printf("Parsed message: %#v", msg.Header)
		if msg.Header.Identifier == lora.PUSH_DATA {
			//err := WriteData(*msg)
			//if err != nil {
			//	log.Print(err.Error())
			//}
			err = msg.Ack()
			if err != nil {
				log.Printf("Error sending ACK: %s", err.Error())
			}
		}
	}
}
