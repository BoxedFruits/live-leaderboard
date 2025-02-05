// Todo:
// - Ping/Pong; Drop client if no response after x seconds
// Auth?
package main

import (
	"log"
	"main/clients"
	"main/matches"
	"main/players"
	"math/rand"
	"net/http"
)

func main() {
	println("Starting server...")
	// Create a match and keep track
	// For right now, make it manually and pass in reference
	// Probably would want to make another service for looking up matches

	// This is kind of like saying one server per process?
	match := &matches.Match{
		MatchId:            rand.Uint32(),
		MatchStatus:        matches.WAITING_FOR_PLAYERS,
		CurrentPlayerCount: 0,
		MaxPlayerCount:     matches.MAX_PLAYER_COUNT,
		Players:            make(map[*clients.Client]*players.PlayerInfo),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		myHandler(w, req, match)
	})
	println("Server is initialized")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func myHandler(w http.ResponseWriter, req *http.Request, testMatch *matches.Match) {
	println("Got a connection")

	client, conn, err := clients.InitializeClient(w, req)

	if err != nil {
		log.Fatal("Could not upgrade connection to websocket")
	}

	//Find/create a match
	testMatch.ConnectClientToMatch(client)

	defer conn.Close()

	for {
		testMatch.ParseMessages(conn, client)
	}
}
