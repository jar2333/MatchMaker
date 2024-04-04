package game

import (
	"log"

	"github.com/gorilla/websocket"

	"github.com/jar2333/MatchMaker/api"
)

// =============================
// = Game playing function
// =============================

// Play the Game g, given two websockets, each corresponding to a player
func Play(g Game, conn1 *websocket.Conn, conn2 *websocket.Conn) {
	// Get player keys for this game
	p1 := g.P1()
	p2 := g.P2()

	// Game loop until game is finished and winner is found:
	for !g.IsFinished() {
		playTurn(g, conn1, p1)

		if g.IsFinished() {
			break
		}

		playTurn(g, conn2, p2)
	}

}

// =============================
// = Helper functions
// =============================

// Turn logic loop
func playTurn(g Game, conn *websocket.Conn, player_id string) {
	var msg []byte
	for {
		// Read message sent by player
		msg = readMessage(conn)

		// Attempt to parse move from message
		move, err := api.ParseMove(msg)
		if err != nil {
			// Parsing error
			sendError(conn, err)
			continue
		}

		// Attempt to play the player's turn
		err = g.PlayTurn(player_id, move)
		if err != nil {
			// Gameplay error
			sendError(conn, err)
		} else {
			break
		}
	}

	// Send game state to player
	sendState(conn, g.State())
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

func sendError(conn *websocket.Conn, err error) {
	// Not yet implemented
}

func sendState(conn *websocket.Conn, state map[string]interface{}) {
	// Not yet implemented
}
