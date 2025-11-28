package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"poker/internal/data"
)

// Define a custom testServer type which embeds an httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initializes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

// Implement a getPlayerScore() method on our custom testServer type. This
// makes a GET request to a given url path using the test server client, and
// returns the response status code and body.
func (ts *testServer) getPlayerScore(t *testing.T, name string) (int, string) {
	rs, err := ts.Client().Get(ts.URL + fmt.Sprintf("/players/%s", name))
	if err != nil {
		// TODO: Do better error message.
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		// TODO: Do better error message.
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, string(body)
}

// Implement a recordWin() method for sending POST request to the test server.
func (ts *testServer) recordWin(t *testing.T, name string) int {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/players/%s", ts.URL, name), nil)

	res, err := ts.Client().Do(req)
	if err != nil {
		// TODO: Do better error message.
		t.Fatal(err)
	}

	return res.StatusCode
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []data.Player) {
	t.Helper()

	league, err := data.NewLeague(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'",
			body, err)
	}

	return
}

// Requests
func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
