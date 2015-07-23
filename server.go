package main

import (
	"fmt"
	"log"
	"net"
)

func StartUDPServer(port int) {
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	CheckError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Print("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			log.Print("Error: ", err)
		}
	}
}
