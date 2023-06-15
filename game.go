package tournament

const (
	TIE = ""
)

type game interface {
	p1() string
	p2() string

	waitForWinner() string

	isFinished() bool

	play(key string, move string)
}
