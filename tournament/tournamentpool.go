package tournament

// Adds a tournament to the pool of tournaments.
// Starts the countdown to its beginning, or otherwise starts necessary goroutines for tournament
func AddTournament(t *Tournament) string {
	return ""
}

func FetchTournament(tournament_key string) (*Tournament, error) {
	return nil, nil
}

func HasTournamentStarted(tournament_key string) bool {
	return true
}
