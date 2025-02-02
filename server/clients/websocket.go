package clients

import (
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func InitializeClient(w http.ResponseWriter, req *http.Request) (*Client, *websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		println("ERRORED WHEN UPGRADING")
		log.Println(err)
		return nil, nil, err
	}
	// defer conn.Close()

	return &Client{
		ClientId:   rand.Uint(),
		ConnStatus: CONNECTED,
	}, conn, nil
}
