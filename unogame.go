package main

type unogame struct {
	_p1 string
	_p2 string

	_finished chan string
}

func (g unogame) p1() string {
	return g._p1
}

func (g unogame) p2() string {
	return g._p2
}

func (g unogame) finished() chan string {
	return g._finished
}

func makeUnoGame(p1 string, p2 string) game {
	return unogame{
		_p1:       p1,
		_p2:       p2,
		_finished: make(chan string),
	}
}
