package game

import "time"

type MockGame struct {
	p1     string
	p2     string
	winner string

	_is_finished bool
}

func MakeMockGame(p1 string, p2 string, winner string) Game {
	return &MockGame{
		p1:     p1,
		p2:     p2,
		winner: winner,
	}
}

func (g *MockGame) P1() string {
	return g.p1
}
func (g *MockGame) P2() string {
	return g.p2
}

func (g *MockGame) WaitForWinner() string {
	time.Sleep(1 * time.Second)
	g._is_finished = true
	return g.winner
}

func (g *MockGame) IsFinished() bool {
	return g._is_finished
}

func (g *MockGame) Play(key string, move map[string]interface{}) bool {
	return true
}

func (g *MockGame) State() map[string]interface{} {
	return make(map[string]interface{})
}
