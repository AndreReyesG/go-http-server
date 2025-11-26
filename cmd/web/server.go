package main

import "net/http"

// Player stores a name with a number of wins.
type Player struct {
	Name string
	Wins int
}

// PlayerStore stores score information about players.
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

// PlayerServer is a HTTP interface for player information.
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("GET /league", http.HandlerFunc(p.leagueHandler))
	router.Handle("GET /players/{player}", http.HandlerFunc(p.showScore))
	router.Handle("POST /players/{player}", http.HandlerFunc(p.processWin))

	p.Handler = router

	return p
}
