package services

import (
	"context"

	"nhooyr.io/websocket"
)

func BuildService[T Light | HomeAssistant](conn *websocket.Conn, ctx context.Context) *T {
	return &T{conn: conn, ctx: ctx}
}
