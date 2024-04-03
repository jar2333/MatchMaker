package game

const (
	TIE = ""
)

type Game interface {
	// Getters for player ids
	P1() string
	P2() string

	// Blocks until winner is determined, returns winner
	// Note: Make optional to denote ties?
	WaitForWinner() string

	// Queries to ask if the game is finished
	IsFinished() bool

	// Attempts to play turn a given move dictionary, returns false on failure to realize play
	Play(player_id string, move map[string]interface{}) bool

	// Returns a state dictionary to send to players
	State() map[string]interface{}
}
