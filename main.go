package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var TOURNAMENT_DATE time.Time

var KEYS = map[string]bool{}

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

	// Load player keys and tournament date from files
	loadKeys()
	loadDate()

	// Create tournament
	tournament := makeTournament()

	// Set http handler and start server
	port := ":" + arguments[1]
	srv := startServer(port, tournament)

	// Wait for tournament to start
	waitForTournamentDate()

	// Play tournament until finishing
	has_tournament_started = true
	tournament.Start()

	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func startServer(port string, t *tournament) *http.Server {
	srv := &http.Server{Addr: port}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(t, w, r)
	})

	go func() {
		fmt.Printf("server started\n")
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
			os.Exit(1)
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			panic(err)
		}
	}()

	return srv
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

func waitForTournamentDate() {
	fmt.Printf("Waiting for tournament to start...\n")
	duration := time.Until(TOURNAMENT_DATE)
	time.Sleep(duration)
}

func loadDate() {
	dat, err := os.ReadFile("./date.txt")
	if err != nil {
		panic(err)
	}

	// Calling Parse() method with its parameters
	tm, e := time.Parse(time.RFC822, string(dat))
	if e != nil {
		panic(e)
	}

	TOURNAMENT_DATE = tm
}

func loadKeys() {
	dat, err := os.ReadFile("./keys.txt")
	if err != nil {
		panic(err)
	}

	lines := splitLines(string(dat))

	for _, l := range lines {
		KEYS[l] = true
	}

}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
