package main

import (
	"log"
	"net/http"
	"os"

	"poker/internal/store"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := store.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating the file system player store, %v", err)
	}

	server := NewPlayerServer(store)

	log.Print("starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
