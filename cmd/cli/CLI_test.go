package main

import (
	"testing"

	"poker/internal/testutils"
)

func TestCLI(t *testing.T) {
	playerStore := &testutils.StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.WinCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}
}
