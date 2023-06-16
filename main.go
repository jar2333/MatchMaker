package tournament

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

// =============================
// = Websocket server
// =============================

func main() {
	defer reg.close()

	// Get command-line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// Load player keys from file
	loadKeys()

	// Set http handler and start server
	http.HandleFunc("/", handler)

	port := ":" + arguments[1]
	err := http.ListenAndServe(port, nil)

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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(games)
		go playGames() // Ends when games channel is closed
		matchMake()    // Ends when tournament is over
	}()

	wg.Wait() // Wait until tournament is over
}

func handler(w http.ResponseWriter, r *http.Request) {
	// If tournament has already started, accept no new connections
	if has_tournament_started {
		return
	}

	// Create websocket connection with client
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// If key was obtained succesfully, register player.
	key, ok := getKey(conn)
	if !ok {
		conn.Close()
		return
	}

	reg.registerPlayer(key, conn)
}

func getKey(conn *websocket.Conn) (string, bool) {
	// Read message from websocket connection
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return "", false
	}

	// Check if message string is a key
	msg := string(p)
	_, ok := KEYS[msg]
	if ok {
		// Write to connection to signal that client has been registered
		if err := conn.WriteMessage(messageType, []byte(REGISTERED_MESSAGE)); err != nil {
			log.Println(err)
			return "", false
		}

		// Return player key
		return msg, true
	}

	return "", false
}

func waitForTournamentStart() {
	// UNIMPLEMENTED
}

func loadKeys() {
	dat, err := os.ReadFile("/keys.txt")
	if err != nil {
		panic(err)
	}

	lines := SplitLines(string(dat))

	for _, l := range lines {
		KEYS[l] = true
	}

}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
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
			played = g.Play(p1, parseMove(msg))
		}
		sendState(conn1, g)

		if g.IsFinished() {
			break
		}

		// Parse player 2's move, perform it, send game state
		played = false
		for !played {
			msg = readMessage(conn2)
			played = g.Play(p2, parseMove(msg))
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
	winner := game.WaitForWinner()

	return winner
}
