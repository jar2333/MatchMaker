package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var reg registry = makeRegistry()
var games chan game = make(chan game)

var tournament_started bool = false
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

	// Start server
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// Accept new connections
	for !tournament_started {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}

	// Play tournament until finishing
	go matchMake()
	go playGames()
	<-tournament_finished
}

func handleConnection(c net.Conn) {
	// Register player to registry, save connection to registry.

	// c.Close()
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
	// p1 := g.p1()
	// p2 := g.p2()

	// Grab a reference to the websockets corresponding to player 1 and player 2

	// Loop until game is finished:
	// 	  Parse player 1's move, perform it, send game state
	// 	  Parse player 2's move, perform it, send game state

	// Mark game as finished, send winning player key (or tie)
	g.finished() <- TIE
}

// =============================
// = Match-making goroutines
// =============================

func matchMake() {
	// Sleep until Tournament start

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
