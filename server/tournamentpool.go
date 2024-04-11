package server

import "github.com/jar2333/MatchMaker/tournament"

// Adds a tournament to the pool of tournaments.
// Starts the countdown to its beginning, or otherwise starts necessary goroutines for tournament
func AddTournament(t *tournament.Tournament) string {
	return ""
}

func FetchTournament(tournament_key string) (*tournament.Tournament, error) {
	return nil, nil
}

func HasTournamentStarted(tournament_key string) bool {
	return true
}
