package tournament

type unogame struct {
	_p1 string
	_p2 string

	_chan        chan string
	_is_finished bool
}

// =========================
// = Generalizable methods
// =========================

func makeUnoGame(p1 string, p2 string) game {
	return unogame{
		_p1:          p1,
		_p2:          p2,
		_chan:        make(chan string),
		_is_finished: false,
	}
}

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
	default:
		return false
	}

	return true
}

func (g *unogame) setWinner(key string) {
	ch := g._chan
	ch <- key
	g._is_finished = true
}

// =========================
// = UNO-SPECIFIC METHODS
// =========================

func (g *unogame) draw(key string) {

}
