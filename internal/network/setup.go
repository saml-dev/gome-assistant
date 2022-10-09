package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"nhooyr.io/websocket"
)

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

func SendAuthMessage(conn websocket.Conn, ctx context.Context) error {
	token := os.Getenv("AUTH_TOKEN")
	msgJson, err := json.Marshal(AuthMessage{MsgType: "auth", AccessToken: token})
	if err != nil {
		return err
	}
	err = conn.Write(ctx, websocket.MessageText, msgJson)
	if err != nil {
		return err
	}
	return nil
}

func WriteMessage[T any](msg T, conn websocket.Conn, ctx context.Context) error {
	msgJson, err := json.Marshal(msg)
	fmt.Println(string(msgJson))
	if err != nil {
		return err
	}

	err = conn.Write(ctx, websocket.MessageText, msgJson)
	if err != nil {
		return err
	}

	return nil
}

func ReadMessage(conn websocket.Conn, ctx context.Context) (string, error) {
	_, msg, err := conn.Read(ctx)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
