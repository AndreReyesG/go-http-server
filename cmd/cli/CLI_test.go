package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"poker/internal/testutils"
)

type GameSpy struct {
	StartCalled bool
	StartedWith int

	FinishedCalled bool
	FinishedWith   string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedCalled = true
	g.FinishedWith = winner
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with Moka as the winner", func(t *testing.T) {
		in := userSends("3", "Moka wins")
		stdout := &bytes.Buffer{}
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertGameFinishedWith(t, game, "Moka")
	})

	t.Run("start game with 9 players and finish game with Milky as the winner", func(t *testing.T) {
		in := userSends("9", "Milky wins")
		game := &GameSpy{}

		cli := NewCLI(in, testutils.DummyStdOut, game)
		cli.PlayPoker()

		assertGameStartedWith(t, game, 9)
		assertGameFinishedWith(t, game, "Milky")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("Hola")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, PlayerPrompt, BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when the winner is declared incorrectly", func(t *testing.T) {
		in := userSends("5", "Moka is a cat")
		stdout := &bytes.Buffer{}
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, PlayerPrompt, BadWinnerInputErrMsg)
	})
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()
	if game.StartedWith != numberOfPlayersWanted {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartedWith)
	}
}

func assertGameFinishedWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	if game.FinishedWith != winner {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishedWith)
	}
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Error("game should not have started")
	}
}

func assertGameNotFinished(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.FinishedCalled {
		t.Error("game should not have finished")
	}
}
