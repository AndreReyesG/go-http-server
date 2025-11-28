package store

import (
	"testing"

	"poker/internal/data"
	"poker/internal/testutils"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, `[
			{"Name": "Moka", "Wins": 10},
			{"Name": "Milky", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		testutils.AssertNoError(t, err)

		got := store.GetLeague()

		want := []data.Player{
			{"Milky", 33},
			{"Moka", 10},
		}

		testutils.AssertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		testutils.AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, `[
			{"Name": "Moka", "Wins": 10},
			{"Name": "Milky", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		testutils.AssertNoError(t, err)

		got := store.GetPlayerScore("Milky")
		want := 33
		testutils.AssertScoreEquals(t, got, want)

	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, `[
			{"Name": "Moka", "Wins": 10},
			{"Name": "Milky", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		testutils.AssertNoError(t, err)

		store.RecordWin("Milky")

		got := store.GetPlayerScore("Milky")
		want := 34
		testutils.AssertScoreEquals(t, got, want)

	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, `[
			{"Name": "Moka", "Wins": 10},
			{"Name": "Milky", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		testutils.AssertNoError(t, err)

		store.RecordWin("Rorro")

		got := store.GetPlayerScore("Rorro")
		want := 1
		testutils.AssertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		testutils.AssertNoError(t, err)
	})
}
