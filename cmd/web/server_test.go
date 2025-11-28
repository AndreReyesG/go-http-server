package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"poker/internal/data"
	"poker/internal/testutils"
)

func TestGETPlayers(t *testing.T) {
	store := testutils.StubPlayerStore{
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
	server := newTestServer(t, NewPlayerServer(&store))
	defer server.Close()

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Moka"

		statusCode := server.recordWin(t, player)
		testutils.AssertStatus(t, statusCode, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin; want 1", len(store.WinCalls))
		}

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner; got: %q; want: %q",
				store.WinCalls[0], player)
		}
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
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertLeague(t, got, wantedLeague)
		testutils.AssertContentType(t, response, jsonContentType)
	})
}
