package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"nhooyr.io/websocket"
)

var ctx, ctxCancel = context.WithTimeout(context.Background(), time.Second*5)
var conn, _, err = websocket.Dial(ctx, "ws://192.168.86.67:8123/api/websocket", nil)

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

func SendAuthMessage() error {
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

func WriteMessage[T any](msg T) error {
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

func ReadMessage() (string, error) {
	_, msg, err := conn.Read(ctx)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
