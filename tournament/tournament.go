package tournament

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jar2333/MatchMaker/game"
)

/**
* ---------------------------------------------------
* A tournament implements its own matchmaking logic.
*----------------------------------------------------
 */
type Tournament struct {
	reg          Registry
	games        chan game.Game
	game_factory func(p1 string, p2 string) game.Game
}

func MakeTournament(game_factory func(p1 string, p2 string) game.Game) *Tournament {
	return &Tournament{
		reg:          MakeRegistry(),
		games:        make(chan game.Game),
		game_factory: game_factory,
	}
}

func (t *Tournament) Register(player_id string, conn *websocket.Conn) {

	t.reg.RegisterPlayer(player_id, conn)
}

func (t *Tournament) Start() {
	defer close(t.games)
	defer t.reg.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Sets up tournament matchmaking, queues up games which need to be played
	go func() {
		defer wg.Done()
		t.matchMake()
	}()

	// Plays games as they are queued to the games channel
	// Ends when games channel is closed
	go t.playGames()

	// Wait until Tournament is over
	wg.Wait()
}

// =============================
// = Game-playing goroutines
// =============================

func (t *Tournament) playGames() {
	for g := range t.games {
		conn1 := t.reg.GetConnection(g.P1())
		conn2 := t.reg.GetConnection(g.P2())
		go game.Play(g, conn1, conn2)
	}
}

// =============================
// = Tournament Match-making goroutines
// =============================

func (t *Tournament) matchMake() {
	// Get Registry
	reg := &t.reg

	// START: create Tournament schedule
	registered := reg.GetRegistered()

	schedule := getSchedule(registered)

	// Play all games
	for _, round := range schedule {
		var wg sync.WaitGroup

		results := make(chan string)

		for _, pair := range round {
			wg.Add(1)
			go func() {
				defer wg.Done()
				winner := t.evalPair(pair)
				results <- winner
			}()
		}

		wg.Wait()

		for winner := range results {
			reg.RecordWin(winner)
		}

	}
}

func (t *Tournament) evalPair(p pair) string {
	p1 := p.p1
	p2 := p.p2

	reg := &t.reg

	if (p1 == EMPTY_KEY || reg.IsDisqualified(p1)) && (p2 == EMPTY_KEY || reg.IsDisqualified(p2)) {
		// No winner
		return EMPTY_KEY
	} else if p1 == EMPTY_KEY || reg.IsDisqualified(p1) {
		// Player 1 loses, Player 2 wins
		return p2
	} else if p2 == EMPTY_KEY || reg.IsDisqualified(p2) {
		// Player 2 loses, Player 1 wins
		return p1
	}

	// Create new game
	game := t.game_factory(p1, p2)

	// Send game to games channel
	t.games <- game

	// Get winner of game from channel
	winner := game.WaitForWinner()

	return winner
}
