package tournament

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var KEYS = map[string]bool{"testkey": true}

var has_tournament_started bool = false

// =============================
// = Websocket server
// =============================

func main() {
	// Get command-line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// Load player keys from file
	loadKeys()

	// Create tournament
	tournament := makeTournament()

	// Set http handler and start server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(tournament, w, r)
	})

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
	tournament.Start()
}

func handler(t *tournament, w http.ResponseWriter, r *http.Request) {
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

	t.Register(key, conn)
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
