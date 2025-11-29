package main

import (
	"log"
	"net/http"

	"poker/internal/store"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := store.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := NewPlayerServer(store)

	log.Print("starting server on :5000")

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000, %v", err)
	}
}
