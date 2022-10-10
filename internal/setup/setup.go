package setup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func SetupConnection(connString string) (*websocket.Conn, context.Context, context.CancelFunc, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)

	// Init websocket connection
	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/api/websocket", connString), nil)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect to websocket at ws://%s/api/websocket. Check IP address and port\n", connString)
		ctxCancel()
		return nil, nil, nil, err
	}

	// Read auth_required message
	_, err = ReadMessage(conn, ctx)
	if err != nil {
		ctxCancel()
		return nil, nil, nil, err
	}

	// Send auth message
	err = SendAuthMessage(conn, ctx)
	if err != nil {
		ctxCancel()
		return nil, nil, nil, err
	}

	// Verify auth message
	err = VerifyAuthResponse(conn, ctx)
	if err != nil {
		fmt.Println("ERROR: Auth token is invalid. Please double check it or create a new token in your Home Assistant profile")
		ctxCancel()
		return nil, nil, nil, err
	}

	return conn, ctx, ctxCancel, err
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
