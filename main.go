package main

import (
	"log"
)

func main() {
	log.Print("Croft is ALIVE")
	go StartUDPServer(1700)
	ServeHTTPOverview()
}
