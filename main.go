package tournament

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var KEYS = map[string]bool{"testkey": true}

var reg registry = makeRegistry()
var games chan game = make(chan game)

var has_tournament_started bool = false
var tournament_finished chan bool = make(chan bool)

// =============================
// = Websocket server
// =============================

func main() {
	defer close(games)
	defer reg.close()

	// Get command-line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// Set http handler and start server
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

	// Wait for tournament to start
	waitForTournamentStart()

	// Play tournament until finishing
	has_tournament_started = true
	go matchMake()
	go playGames()
	<-tournament_finished

}

func handler(w http.ResponseWriter, r *http.Request) {
	// If tournament has already started, accept no new connections
	if has_tournament_started {
		return
	}

	// Create websocket connection with client
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Error when creating connection, close
		log.Println(err)
		conn.Close()
		return
	}

	key := getKey(conn)

	reg.registerPlayer(key, conn)
}

func getKey(conn *websocket.Conn) string {
	for {
		// Read message from websocket connection
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		// Check if message string is a key
		msg := string(p)
		_, ok := KEYS[msg]
		if ok {
			// Write to connection to signal that client has been registered
			if err := conn.WriteMessage(messageType, []byte(REGISTERED_MESSAGE)); err != nil {
				log.Println(err)
				continue
			}

			// Return player key
			return msg
		}

	}
}

func waitForTournamentStart() {
	// UNIMPLEMENTED
}

// =============================
// = Game-playing goroutines
// =============================

func playGames() {
	for g := range games {
		go playGame(g)
	}
}

func playGame(g game) {
	// Get player keys for this game
	p1 := g.p1()
	p2 := g.p2()

	// Grab a reference to the websockets corresponding to player 1 and player 2
	conn1 := reg.getConnection(p1)
	conn2 := reg.getConnection(p2)

	// Loop until game is finished and winner is found:
	var move string
	var winner string = TIE
	for !g.is_finished() {
		// Parse player 1's move, perform it, send game state
		move = readMove(conn1)
		g.play(p1, move)
		sendState(conn1, g)

		if g.is_finished() {
			break
		}

		// Parse player 2's move, perform it, send game state
		move = readMove(conn2)
		g.play(p2, move)
		sendState(conn2, g)
	}

	// Mark game as finished, send winning player key (or tie)
	// NOTE: MAYBE SHOULD BE DONE _INSIDE_ THE GAME METHODS?
	g.finished() <- winner
}

func readMove(conn *websocket.Conn) string {
	// Not yet implemented
	return ""
}

func sendState(conn *websocket.Conn, g game) {
	// Not yet implemented
}

// =============================
// = Match-making goroutines
// =============================

func matchMake() {
	// START: create tournament schedule
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
				winner := evalGame(pair)
				results <- winner
			}()
		}

		wg.Wait()

		for winner := range results {
			reg.recordWin(winner)
		}

	}

	// All games played!
	tournament_finished <- true
}

func evalGame(p pair) string {
	p1 := p.p1
	p2 := p.p2

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
	game := makeGame(p1, p2)

	// Send game to games channel
	games <- game

	// Get winner of game from channel
	var ch chan string = game.finished()
	var winner string = <-ch

	return winner
}
