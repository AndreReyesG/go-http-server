package main

import (
	"poker/internal/data"
)

type CLI struct {
	store data.PlayerStore
}

func (cli *CLI) PlayPoker() {
	cli.store.RecordWin("Milky")
}
