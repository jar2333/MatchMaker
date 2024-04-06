package matchmaker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jar2333/MatchMaker/tournament"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

	// Set http handler and start server
	port := ":" + arguments[1]
	srv := startServer(port)

	// Shut server down
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func startServer(port string) *http.Server {
	srv := &http.Server{Addr: port}

	http.HandleFunc("/", handler)

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

// ====================
// == HTTP Handlers
// ====================
func handler(w http.ResponseWriter, r *http.Request) {
	// Add authentication here!

	// Check if player requested tournament join, or regular matchmaking (URI API).
	// Check what game the player has requested as well, to get adequate game factory
	// How can this be added through an API?
	// Work on this.
	request_type := "tournament_register"

	// Create websocket connection with client
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	switch request_type {
	case "match_make":
		{

		}
	case "create_tournament":
		{

		}
	case "register_tournament":
		{
			// If key was obtained succesfully, register player to tournament.
			key, ok := getTournamentKey(conn)
			if !ok {
				conn.Close()
				return
			}

			// Attempt to find a tournament with the given key
			t, err := tournament.FetchTournament(key)
			if err != nil {
				conn.Close()
				return
			}

			// If tournament has already started, accept no new connections
			if tournament.HasTournamentStarted(key) {
				conn.Close()
				return
			}

			// Get a unique identifier for the player, and associate it with the websocket connection
			player_id := r.Host
			t.Register(player_id, conn)
		}
	default:
		{
			// No appropriate procedures available
			conn.Close()
		}
	}

}

// ==================
// == Helpers
// ==================
func getTournamentKey(conn *websocket.Conn) (string, bool) {
	// Read message from websocket connection
	// 	messageType, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return "", false
	// 	}

	// 	// Check if message string is a key
	// 	msg := string(p)
	// 	_, ok := KEYS[msg]
	// 	if ok {
	// 		// Write to connection to signal that client has been registered
	// 		if err := conn.WriteMessage(messageType, []byte(api.REGISTERED_MESSAGE)); err != nil {
	// 			log.Println(err)
	// 			return "", false
	// 		}

	// 		// Return player key
	// 		return msg, true
	// 	}

	return "", false
}

// func waitForTournamentDate() {
// 	fmt.Printf("Waiting for tournament to start...\n")
// 	duration := time.Until(TOURNAMENT_DATE)
// 	time.Sleep(duration)
// }
