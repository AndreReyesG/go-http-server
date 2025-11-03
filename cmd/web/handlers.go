package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonContentType = "application/json"

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := r.PathValue("player")
	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := r.PathValue("player")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
