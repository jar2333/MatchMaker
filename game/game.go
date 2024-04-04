package game

type GameError struct {
	Input string
}

func (e *GameError) Error() string {
	return "Invalid move: " + e.Input
}

type Game interface {
	// Getters for player ids
	P1() string
	P2() string

	// Blocks until winner is determined, returns winner
	// Note: Make optional to denote ties.
	WaitForWinner() string

	// Queries to ask if the game is finished
	IsFinished() bool

	// Attempts to play the given player's turn a given move dictionary
	// Returns an error type
	PlayTurn(player_id string, move map[string]interface{}) error

	// Returns a state dictionary to send to players
	State() map[string]interface{}
}
