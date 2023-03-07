package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type BaseMessage struct {
	Type string `json:"type"`
	Id   int64  `json:"id"`
}

type ChanMsg struct {
	Id   int64
	Type string
	Raw  []byte
}

func ListenWebsocket(conn *websocket.Conn, ctx context.Context, c chan ChanMsg) {
	for {
		bytes, err := ReadMessage(conn, ctx)

		if err != nil {
			log.Default().Println("Error reading from websocket:", err)
			close(c)
			break
		}

		base := BaseMessage{}
		json.Unmarshal(bytes, &base)
		chanMsg := ChanMsg{
			Type: base.Type,
			Id:   base.Id,
			Raw:  bytes,
		}

		c <- chanMsg
	}
}
