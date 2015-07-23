package main

import (
	"log"
)

func main() {
	log.Print("Croft is ALIVE")
	StartUDPServer(1700)
	ServeHTTPOverview()
}
