package game

type MockGame struct {
	p1     string
	p2     string
	winner string
}

func MakeMockGame(p1 string, p2 string, winner string) Game {
	return &MockGame{
		p1:     p1,
		p2:     p2,
		winner: winner,
	}
}

func (g *MockGame) P1() string {
	return ""
}
func (g *MockGame) P2() string {
	return ""
}

func (g *MockGame) WaitForWinner() string {
	return ""
}

func (g *MockGame) IsFinished() bool {
	return false
}

func (g *MockGame) Play(key string, move map[string]interface{}) bool {
	return true
}

func (g *MockGame) State() map[string]interface{} {
	return make(map[string]interface{})
}
