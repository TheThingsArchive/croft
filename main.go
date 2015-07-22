package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func main() {
	log.Print("Croft is ALIVE")

	ServerAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:1700")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	log.Printf("%#v", ServerConn)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	go func() {
		for {
			n, addr, err := ServerConn.ReadFromUDP(buf)
			log.Print("Received ", string(buf[0:n]), " from ", addr)

			if err != nil {
				log.Print("Error: ", err)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Croft")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
