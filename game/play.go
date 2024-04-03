package game

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/jar2333/MatchMaker/api"
)

func Play(g Game, conn1 *websocket.Conn, conn2 *websocket.Conn) {
	// Get player keys for this game
	p1 := g.P1()
	p2 := g.P2()

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

func sendState(conn *websocket.Conn, g Game) {
	// Not yet implemented
}
