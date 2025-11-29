package main

import (
	"strings"
	"testing"

	"poker/internal/testutils"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Moka wins\n")
	playerStore := &testutils.StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	testutils.AssertPlayerWin(t, playerStore, "Moka")
}
