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

	// Plays a given move dictionary
	Play(key string, move map[string]interface{}) bool

	// Returns a state dictionary
	State() map[string]interface{}
}
