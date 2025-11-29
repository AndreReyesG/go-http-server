package main

import (
	"bufio"
	"io"
	"strings"

	"poker/internal/data"
)

type CLI struct {
	store data.PlayerStore
	in    io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.store.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
