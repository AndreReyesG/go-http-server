package main

import (
	"bytes"
	"strings"
	"testing"

	"poker/internal/testutils"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &testutils.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartCalled bool
	StartedWith int

	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	t.Run("record Moka win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nMoka wins\n")
		game := &GameSpy{}

		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Moka" {
			t.Errorf("expected finish called with 'Moka' but got %q", game.FinishedWith)
		}
	})

	t.Run("record Milky win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nMilky wins\n")
		game := &GameSpy{}

		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Milky" {
			t.Errorf("expected finish called with 'Milky' but got %q", game.FinishedWith)
		}
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Hola\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Error("game should not have started")
		}

		gotPrompt := stdout.String()
		wantPrompt := PlayerPrompt + "you're so silly"

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}
	})
}
