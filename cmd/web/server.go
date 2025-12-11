package main

import (
	"fmt"
	"html/template"
	"net/http"

	"poker/internal/data"
)

// PlayerServer is a HTTP interface for player information.
type PlayerServer struct {
	store data.PlayerStore
	http.Handler
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store data.PlayerStore) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s, %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("GET /league", http.HandlerFunc(p.leagueHandler))
	router.Handle("GET /players/{player}", http.HandlerFunc(p.showScore))
	router.Handle("POST /players/{player}", http.HandlerFunc(p.processWin))
	router.Handle("GET /game", http.HandlerFunc(p.game))

	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}
