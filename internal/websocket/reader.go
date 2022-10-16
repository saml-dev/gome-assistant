package websocket

import (
	"context"

	"nhooyr.io/websocket"
)

type BaseMessage struct {
	MsgType string `json:"type"`
	Other   any
}

func ReadWebsocket(ws *websocket.Conn, ctx context.Context) {

}
