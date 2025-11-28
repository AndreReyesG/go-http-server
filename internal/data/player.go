package data

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
