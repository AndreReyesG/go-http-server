package main

import (
	"strings"
	"testing"

	"poker/internal/testutils"
)

func TestCLI(t *testing.T) {
	t.Run("record Moka win from user input", func(t *testing.T) {
		in := strings.NewReader("Moka wins\n")
		playerStore := &testutils.StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		testutils.AssertPlayerWin(t, playerStore, "Moka")
	})

	t.Run("record Milky win from user input", func(t *testing.T) {
		in := strings.NewReader("Milky wins\n")
		playerStore := &testutils.StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		testutils.AssertPlayerWin(t, playerStore, "Milky")
	})
}
