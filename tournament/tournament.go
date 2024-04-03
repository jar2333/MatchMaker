package tournament

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jar2333/MatchMaker/api"
	"github.com/jar2333/MatchMaker/game"
)

type Tournament struct {
	reg          registry
	games        chan game.Game
	game_factory func(p1 string, p2 string) game.Game
}

func MakeTournament(game_factory func(p1 string, p2 string) game.Game) *Tournament {
	return &Tournament{
		reg:          makeRegistry(),
		games:        make(chan game.Game),
		game_factory: game_factory,
	}
}

func (t *Tournament) Register(key string, conn *websocket.Conn) {
	t.reg.registerPlayer(key, conn)
}

func (t *Tournament) Start() {
	defer close(t.games)
	defer t.reg.close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		go t.playGames() // Ends when games channel is closed
		t.matchMake()    // Ends when Tournament is over
	}()

	wg.Wait() // Wait until Tournament is over
}

// =============================
// = Game-playing goroutines
// =============================

func (t *Tournament) playGames() {
	for g := range t.games {
		go t.playGame(g)
	}
}

func (t *Tournament) playGame(g game.Game) {
	// Get registry
	reg := &t.reg

	// Get player keys for this game
	p1 := g.P1()
	p2 := g.P2()

	// Grab a reference to the websockets corresponding to player 1 and player 2
	conn1 := reg.getConnection(p1)
	conn2 := reg.getConnection(p2)

	// Game loop until game is finished and winner is found:
	var msg []byte
	var played bool
	for !g.IsFinished() {
		// Parse player 1's move, perform it, send game state
		played = false
		for !played {
			msg = readMessage(conn1)
			played = g.Play(p1, api.ParseMove(msg))
		}
		sendState(conn1, g)

		if g.IsFinished() {
			break
		}

		// Parse player 2's move, perform it, send game state
		played = false
		for !played {
			msg = readMessage(conn2)
			played = g.Play(p2, api.ParseMove(msg))
		}
		sendState(conn2, g)
	}

}

func readMessage(conn *websocket.Conn) []byte {
	for {
		// Read message from websocket connection
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}
		return msg
	}
}

func sendState(conn *websocket.Conn, g game.Game) {
	// Not yet implemented
}

// =============================
// = Match-making goroutines
// =============================

func (t *Tournament) matchMake() {
	// Get registry
	reg := &t.reg

	// START: create Tournament schedule
	registered := reg.getRegistered()

	schedule := getSchedule(registered)

	// Play all games
	for _, round := range schedule {
		var wg sync.WaitGroup

		results := make(chan string)

		for _, pair := range round {
			wg.Add(1)
			go func() {
				defer wg.Done()
				winner := t.evalGame(pair)
				results <- winner
			}()
		}

		wg.Wait()

		for winner := range results {
			reg.recordWin(winner)
		}

	}
}

func (t *Tournament) evalGame(p pair) string {
	p1 := p.p1
	p2 := p.p2

	reg := &t.reg

	if (p1 == EMPTY_KEY || reg.isDisqualified(p1)) && (p2 == EMPTY_KEY || reg.isDisqualified(p2)) {
		// No winner
		return EMPTY_KEY
	} else if p1 == EMPTY_KEY || reg.isDisqualified(p1) {
		// Player 1 loses, Player 2 wins
		return p2
	} else if p2 == EMPTY_KEY || reg.isDisqualified(p2) {
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
