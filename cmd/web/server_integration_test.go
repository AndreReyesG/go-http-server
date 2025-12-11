package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"poker/internal/data"
	"poker/internal/store"
	"poker/internal/testutils"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := testutils.CreateTempFile(t, "[]")
	defer cleanDatabase()
	playerStore, err := store.NewFileSystemPlayerStore(database)

	testutils.AssertNoError(t, err)

	server := mustMakePlayerServer(t, playerStore)
	player := "Moka"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		testutils.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []data.Player{
			{"Moka", 3},
		}
		testutils.AssertLeague(t, got, want)
	})
}
