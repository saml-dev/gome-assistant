// Package websocket is used to interact with the Home Assistant
// websocket API. All HA interaction is done via websocket
// except for cases explicitly called out in http package
// documentation.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"nhooyr.io/websocket"
)

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

func WriteMessage[T any](msg T, conn *websocket.Conn, ctx context.Context) error {
	msgJson, err := json.Marshal(msg)
	// fmt.Println(string(msgJson))
	if err != nil {
		return err
	}

	err = conn.Write(ctx, websocket.MessageText, msgJson)
	if err != nil {
		return err
	}

	return nil
}

func ReadMessage(conn *websocket.Conn, ctx context.Context) (string, error) {
	_, msg, err := conn.Read(ctx)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func SetupConnection(connString string) (*websocket.Conn, context.Context, context.CancelFunc) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3)

	// Init websocket connection
	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/api/websocket", connString), nil)
	if err != nil {
		ctxCancel()
		log.Fatalf("ERROR: Failed to connect to websocket at ws://%s/api/websocket. Check IP address and port\n", connString)
	}

	// Read auth_required message
	_, err = ReadMessage(conn, ctx)
	if err != nil {
		ctxCancel()
		log.Fatalln("Unknown error creating websocket client")
	}

	// Send auth message
	err = SendAuthMessage(conn, ctx)
	if err != nil {
		ctxCancel()
		log.Fatalln("Unknown error creating websocket client")
	}

	// Verify auth message
	err = VerifyAuthResponse(conn, ctx)
	if err != nil {
		ctxCancel()
		log.Fatalln("ERROR: Auth token is invalid. Please double check it or create a new token in your Home Assistant profile")
	}

	return conn, ctx, ctxCancel
}

func SendAuthMessage(conn *websocket.Conn, ctx context.Context) error {
	token := os.Getenv("AUTH_TOKEN")
	err := WriteMessage(AuthMessage{MsgType: "auth", AccessToken: token}, conn, ctx)
	if err != nil {
		return err
	}
	return nil
}

type authResponse struct {
	MsgType string `json:"type"`
	Message string `json:"message"`
}

func VerifyAuthResponse(conn *websocket.Conn, ctx context.Context) error {
	_, msg, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	var authResp authResponse
	json.Unmarshal(msg, &authResp)
	if authResp.MsgType != "auth_ok" {
		return errors.New("invalid auth token")
	}

	return nil
}
