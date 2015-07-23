package main

import (
	"log"
)

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
