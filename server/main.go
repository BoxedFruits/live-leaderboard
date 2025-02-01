// Todo:
// - Ping/Pong; Drop client if no response after x seconds
// Auth?
package main

import (
	"fmt"
	"log"
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

type ClientConnStatuses string

const (
	CONNETCING    ClientConnStatuses = "CONNECTING"
	CONNECTED     ClientConnStatuses = "CONNECTED"
	DISCONNECTING ClientConnStatuses = "DISCONNECTING"
)

type Client struct {
	ClientId   uint
	ConnStatus ClientConnStatuses
}

type Leaderboard struct {
	PlayerInfos []PlayerInfo
}

type PlayerInfo struct {
	Client Client
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

const MAX_PLAYER_COUNT = 4

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	//otherwise, keep track of client
	client := Client{
		ClientId:   rand.Uint(),
		ConnStatus: CONNECTED,
	}

	//Find/create a match
	testMatch.ConnectClientToMatch(client)

	defer conn.Close()

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Write message to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			break
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

func (match *Match) ConnectClientToMatch(c Client) {
	if match.CurrentPlayerCount == MAX_PLAYER_COUNT {
		fmt.Println("Too many players. Cannot add another to this server")
		return
	}

	p := PlayerInfo{
		Client: c,
		Kills:  0,
		Deaths: 0,
	}

	println("Client conencted: ", c.ClientId)
	match.Players = append(match.Players, p)
}
