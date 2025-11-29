package main

import (
	"fmt"
	"log"
	"os"

	"poker/internal/store"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := store.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating the file system player store, %v", err)
	}

	game := CLI{store, os.Stdin}
	game.PlayPoker()
}
