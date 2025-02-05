package matches

import (
	"encoding/json"
	"fmt"
	"log"
	"main/clients"
	"main/players"
	"strconv"

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

type Match struct {
	MatchId            uint32
	MatchStatus        MatchStatuses
	CurrentPlayerCount uint8
	MaxPlayerCount     uint8
	Players            map[*clients.Client]*players.PlayerInfo
}

type Leaderboard struct {
	PlayerInfos []players.PlayerInfo
}

type WsRequest struct {
	Command string      `json:"command"`
	Value   interface{} `json:"value"`
}

const MAX_PLAYER_COUNT = 4

func (match *Match) ParseMessages(conn *websocket.Conn, client *clients.Client) {
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
	case "IncrementKillCount":
		{
			println("Incrementing kill count")
			match.Players[client].Kills += 1
		}
	case "GetLeaderboard":
		{
			err = conn.WriteMessage(1, match.GetLeaderboard())
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

// Returns map of clientId: playerInfo, Kills,Deaths
func (match *Match) GetLeaderboard() []byte {
	leaderboard := make(map[string]*players.PlayerInfo)

	for k, v := range match.Players {
		leaderboard[strconv.FormatUint(uint64(k.ClientId), 10)] = v
	}

	json, err := json.Marshal(leaderboard)

	if err != nil {
		println("Error making leaderboard", err)
	}

	return json
}

// func initMatch() *Match {
// 	return &Match{
// 		MatchId:            rand.Uint32(),
// 		MatchStatus:        WAITING_FOR_PLAYERS,
// 		CurrentPlayerCount: 0,
// 		MaxPlayerCount:     MAX_PLAYER_COUNT,
// 		Players:            make(map[*clients.Client]*players.PlayerInfo),
// 	}
// }

func (match *Match) ConnectClientToMatch(c *clients.Client) {
	if match.CurrentPlayerCount == MAX_PLAYER_COUNT {
		fmt.Println("Too many players. Cannot add another to this server")
		return
	}

	p := &players.PlayerInfo{
		Kills:  0,
		Deaths: 0,
	}

	println("Client connected: ", c.ClientId)
	match.Players[c] = p
}
