package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"poker/internal/data"
	"poker/internal/testutils"

	"github.com/gorilla/websocket"
)

func mustMakePlayerServer(t *testing.T, store data.PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatalf("problem creating player server, %v", err)
	}
	return server
}

func TestGETPlayers(t *testing.T) {
	store := testutils.StubPlayerStore{
		map[string]int{
			"Moka":  20,
			"Milky": 10,
		},
		nil,
		nil,
	}
	server := newTestServer(t, mustMakePlayerServer(t, &store))
	defer server.Close()

	// Create a slice of anonymous structs containing the test case name,
	// player name, expected HTTP status code, and expected score.
	tests := []struct {
		name               string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			name:               "returns Moka's score",
			player:             "Moka",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "returns Milky's score",
			player:             "Milky",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "returns 404 on missing players",
			player:             "Rorro",
			expectedHTTPStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := server.getPlayerScore(t, tt.player)

			testutils.AssertStatus(t, statusCode, tt.expectedHTTPStatus)
			testutils.AssertResponseBody(t, body, tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := testutils.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := newTestServer(t, mustMakePlayerServer(t, &store))
	defer server.Close()

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Moka"

		statusCode := server.recordWin(t, player)
		testutils.AssertStatus(t, statusCode, http.StatusAccepted)
		testutils.AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []data.Player{
			{"Milky", 32},
			{"Moka", 12},
			{"Benito", 94},
		}
		store := testutils.StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertLeague(t, got, wantedLeague)
		testutils.AssertContentType(t, response, jsonContentType)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &testutils.StubPlayerStore{})

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &testutils.StubPlayerStore{}
		winner := "Moka"
		server := httptest.NewServer(mustMakePlayerServer(t, store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, winner)

		// NOTE: Puttting arbitary sleeps into tests is very bad practice.
		time.Sleep(10 * time.Millisecond)
		testutils.AssertPlayerWin(t, store, winner)
	})
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
