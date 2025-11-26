package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
		{"Name": "Moka", "Wins": 10},
		{"Name": "Milky", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Moka", 10},
			{"Milky", 33},
		}

		assertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}
