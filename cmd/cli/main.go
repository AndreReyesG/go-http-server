package main

import (
	"fmt"
	"log"
	"os"

	"poker/internal/store"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := store.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {name} wins to record a win")

	game := NewTexasHoldem(BlindAlerterFunc(StdOutAlerter), store)
	cli := NewCLI(os.Stdin, os.Stdout, game)

	cli.PlayPoker()
}
