package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Moka":  20,
			"Milky": 10,
		},
		nil,
		nil,
	}
	server := newTestServer(t, NewPlayerServer(&store))
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

			assertStatus(t, statusCode, tt.expectedHTTPStatus)
			assertResponseBody(t, body, tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := newTestServer(t, NewPlayerServer(&store))
	defer server.Close()

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Moka"

		statusCode := server.recordWin(t, player)
		assertStatus(t, statusCode, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin; want 1", len(store.winCalls))
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner; got: %q; want: %q",
				store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Milky", 32},
			{"Moka", 12},
			{"Benito", 94},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
	})
}
