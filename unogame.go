package tournament

type unogame struct {
	_p1 string
	_p2 string

	_chan        chan string
	_is_finished bool
}

func (g unogame) p1() string {
	return g._p1
}

func (g unogame) p2() string {
	return g._p2
}

func (g unogame) waitForWinner() string {
	ch := g._chan
	winner := <-ch
	return winner
}

func (g unogame) isFinished() bool {
	return g._is_finished
}

func (g unogame) play(key string, move string) {
	// Not yet implemented
}

func makeUnoGame(p1 string, p2 string) game {
	return unogame{
		_p1:          p1,
		_p2:          p2,
		_chan:        make(chan string),
		_is_finished: false,
	}
}
