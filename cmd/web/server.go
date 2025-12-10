package main

import (
	"net/http"

	"poker/internal/data"
)

// PlayerServer is a HTTP interface for player information.
type PlayerServer struct {
	store data.PlayerStore
	http.Handler
}

func NewPlayerServer(store data.PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("GET /league", http.HandlerFunc(p.leagueHandler))
	router.Handle("GET /players/{player}", http.HandlerFunc(p.showScore))
	router.Handle("POST /players/{player}", http.HandlerFunc(p.processWin))
	router.Handle("GET /game", http.HandlerFunc(p.game))

	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p
}
