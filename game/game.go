package game

const (
	TIE = ""
)

type Game interface {
	P1() string
	P2() string

	WaitForWinner() string

	IsFinished() bool

	Play(key string, move map[string]interface{}) bool
}
