package matchmaking

import (
	"github.com/gorilla/websocket"

	"github.com/jar2333/MatchMaker/game"
)

/**
*
* The basic matchmaking stuff, can be put in its own package, to be contrasted with tournament matching (which follows a fixed schedule)
*
 */

type Player struct {
	id   string
	conn *websocket.Conn
}

type MatchMaker struct {
	players      chan Player
	games        chan game.Game
	game_factory func(p1 string, p2 string) game.Game
}

func MakeMatchMaker(game_factory func(p1 string, p2 string) game.Game) *MatchMaker {
	m := MatchMaker{
		players:      make(chan Player),
		games:        make(chan game.Game),
		game_factory: game_factory,
	}
	m.matchMake()
	return &m
}

func (m *MatchMaker) FindGame(player_id string, conn *websocket.Conn) game.Game {
	player := Player{player_id, conn}

	// Add player to channel of players to be matched
	m.players <- player

	// Pop a game from the channel of games that are contantly being created
	game := <-m.games

	return game
}

func (m *MatchMaker) matchMake() {
	// Implements strategy to match players seeking a game
	// EXTENSIBLE LATER
	go func() {
		for {
			p1 := <-m.players
			p2 := <-m.players

			m.games <- m.game_factory(p1.id, p2.id)
		}
	}()
}
