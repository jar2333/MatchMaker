package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// =============================
// = Websocket server
// =============================

func StartServer(port string) *http.Server {
	srv := &http.Server{Addr: port}

	// Register HTTP endpoints

	// Register WebSocket
	http.HandleFunc("/game", getGameConnection)
	// http.HandleFunc("/tournament/{id}/game", getTournamentGameConnection)

	// Start server
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
func getGameConnection(w http.ResponseWriter, r *http.Request) {
	// Add authentication here!

	// Create websocket connection with client
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Matchmake, start game (nonblocking to allow return)

}

// func getTournamentGameConnection(w http.ResponseWriter, r *http.Request) {
// 	// // Add authentication here!

// 	// // Get params
// 	// tournamentID := r.PathValue("id")

// 	// // Create websocket connection with client
// 	// conn, err := upgrader.Upgrade(w, r, nil)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	conn.Close()
// 	// 	return
// 	// }

// 	// // Get tournament from id
// 	// t, err = FetchTournament(tournamentID)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	conn.Close()
// 	// 	return
// 	// }

// 	// // Matchmake, start game (nonblocking to allow return)

// }

// func createTournamentRegistration(w http.ResponseWriter, r *http.Request) {
// 	// Add authentication here!

// 	// If key was obtained succesfully, register player to tournament.

// 	// Attempt to find a tournament with the given key
// 	t, err := tournament.FetchTournament(key)
// 	if err != nil {
// 		// Send error response
// 		log.Println(err)
// 		return
// 	}

// 	// If tournament has already started, accept no new connections
// 	if tournament.HasTournamentStarted(key) {
// 		// Send error response
// 		return
// 	}

// 	// Get a unique identifier for the player, and associate it with the websocket connection
// 	player_id := r.Host
// 	t.Register(player_id, conn)

// }
