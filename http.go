package main

import (
	"fmt"
	"log"
	"net/http"
)

func ServeHTTPOverview() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Croft")
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
