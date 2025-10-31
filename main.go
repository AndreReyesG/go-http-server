package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}
	log.Print("starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
