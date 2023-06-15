package tournament

const (
	TIE = ""
)

type game interface {
	p1() string
	p2() string

	finished() chan string

	is_finished() bool

	play(key string, move string)
}
