package main

import "net/http"

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

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
	router.Handle("POST /players/", http.HandlerFunc(p.processWin))

	p.Handler = router

	return p
}
