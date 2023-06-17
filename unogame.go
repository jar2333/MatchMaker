package main

import gouno "github.com/jar2333/goUno/goUno"

type unogame struct {
	_p1 string
	_p2 string

	_chan        chan string
	_is_finished bool

	gouno.UnoGame
}

func makeUnoGame(p1 string, p2 string) game {
	g := unogame{
		_p1:          p1,
		_p2:          p2,
		_chan:        make(chan string),
		_is_finished: false,
	}

	// Initialize game
	g.Initialize([]string{p1, p2})

	// Mark players as human
	for i := 0; i < 2; i++ {
		g.Players[i].Human = true
	}

	return g
}

// =========================
// = Generalizable methods
// =========================

func (g unogame) P1() string {
	return g._p1
}

func (g unogame) P2() string {
	return g._p2
}

func (g unogame) WaitForWinner() string {
	ch := g._chan
	winner := <-ch
	return winner
}

func (g unogame) IsFinished() bool {
	return g._is_finished
}

func (g unogame) Play(key string, move map[string]interface{}) bool {
	if move == nil {
		return false
	}

	typ, ok := move["type"]
	if !ok {
		return false
	}

	switch typ {
	case "draw":
		g.draw(key)
	case "play":
		if pos, ok := move["position"]; !ok {
			return false
		} else {
			switch pos.(type) {
			case int:
				g.discard(key, pos.(int))
			default:
				return false
			}
		}
	default:
		return false
	}

	return true
}

func (g *unogame) setWinner(key string) {
	g._is_finished = true
	ch := g._chan
	ch <- key
}

// =========================
// = UNO-SPECIFIC METHODS
// =========================

func (g *unogame) draw(key string) {

}

func (g *unogame) discard(key string, position int) {

}
