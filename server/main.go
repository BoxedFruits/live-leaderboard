// Todo:
// - Ping/Pong; Drop client if no response after x seconds
// Auth?
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/clients"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
)

type MatchStatuses string

const (
	STARTING            MatchStatuses = "STARTING_MATCH"
	FINISHED            MatchStatuses = "FINISHED_MATCH"
	WAITING_FOR_PLAYERS MatchStatuses = "WAITING_FOR_PLAYERS"
	IN_PROGRESS         MatchStatuses = "IN_PROGRESS"
	SERVER_BAD          MatchStatuses = "ERROR_IN_SERVER"
)

type Leaderboard struct {
	PlayerInfos []PlayerInfo
}

type PlayerInfo struct {
	Client *clients.Client
	Kills  uint8
	Deaths uint8
}

type PlayerEvents string

const (
	INCREMENT_KILL_COUNT  PlayerEvents = "INCREMENT_KILL_COUNT"
	INCREMENT_DEATH_COUNT PlayerEvents = "INCREMENT_DEATH_COUNT"
)

type Match struct {
	MatchId            uint32
	MatchStatus        MatchStatuses
	CurrentPlayerCount uint8
	MaxPlayerCount     uint8
	Players            []PlayerInfo
	Leaderboard        Leaderboard
}

type WsRequest struct {
	Command string      `json:"command"`
	Value   interface{} `json:"value"`
}

const MAX_PLAYER_COUNT = 4

func main() {
	println("Starting server...")
	// Create a match and keep track
	// For right now, make it manually and pass in reference
	// Probably would want to make another service for looking up matches

	// This is kind of like saying one server per process?
	match := initMatch()

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		myHandler(w, req, match)
	})
	println("Server is initialized")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func myHandler(w http.ResponseWriter, req *http.Request, testMatch *Match) {
	println("Got a connection")

	client, conn, err := clients.InitializeClient(w, req)

	if err != nil {
		log.Fatal("Could not upgrade connection to websocket")
	}

	//Find/create a match
	testMatch.ConnectClientToMatch(client)

	defer conn.Close()

	for {
		parseMessages(conn)
		// Read message from client
		// messageType, message, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }
		// log.Printf("recv: %s", message)
		//Parse message into a command

		// jsonStr, err := json.Marshal(message)
		// if err != nil {
		// fmt.Println(err)
		// return
		// }
		//
		// println(jsonStr)
		//
		// Write message to client
		// err = conn.WriteMessage(messageType, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func parseMessages(conn *websocket.Conn) {

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		// break
	}

	log.Printf("recv: %s", message)

	var jsonStr WsRequest
	errMarshal := json.Unmarshal(message, &jsonStr)

	if errMarshal != nil {
		fmt.Println(errMarshal)
		return
	}
	println(jsonStr.Command)
	switch jsonStr.Command {
	case "command":
		{
			println("in foo")
		}
	case "Hello, Server!":
		{
			println("In the second place")
		}
	}
}

func initMatch() *Match {
	return &Match{
		MatchId:            rand.Uint32(),
		MatchStatus:        WAITING_FOR_PLAYERS,
		CurrentPlayerCount: 0,
		MaxPlayerCount:     MAX_PLAYER_COUNT,
	}
}

func (match *Match) ConnectClientToMatch(c *clients.Client) {
	if match.CurrentPlayerCount == MAX_PLAYER_COUNT {
		fmt.Println("Too many players. Cannot add another to this server")
		return
	}

	p := PlayerInfo{
		Client: c,
		Kills:  0,
		Deaths: 0,
	}

	println("Client connected: ", c.ClientId)
	match.Players = append(match.Players, p)
}

func (match *Match) GetLeaderboard() Leaderboard {
	return match.Leaderboard
}
