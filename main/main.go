package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jar2333/MatchMaker/server"
)

func main() {
	// Get command-line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// Set http handler and start server
	port := ":" + arguments[1]
	srv := server.StartServer(port)

	// Shut server down
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
