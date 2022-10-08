package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saml-dev/gome-assistant/internal/network"
	"nhooyr.io/websocket"
)

var c = context.Background()
var b, x = context.WithTimeout(c, time.Second)

const a = time.April

var ctx, ctxCancel = context.WithTimeout(context.Background(), time.Second*5)
var conn, _, err = websocket.Dial(ctx, "ws://192.168.86.67:8123/api/websocket", nil)

type Light struct {
	EntityId string
}

type LightOnRequest struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Service string `json:"service"`
	Target  struct {
		EntityId string `json:"entity_id"`
	} `json:"target"`
}

func NewLightOnRequest(entity_id string) LightOnRequest {
	req := LightOnRequest{
		Id:      5,
		Type:    "call_service",
		Domain:  "light",
		Service: "turn_on",
	}
	req.Target.EntityId = entity_id
	return req
}

func (l *Light) TurnOn() error {
	// req := json.Marshal()
	return nil
}

func main() {

	sched := Schedule{
		RunEvery: Daily,
	}
	defer ctxCancel()
	if err != nil {
		panic(err)
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")
	// _, _, err = c.Reader(ctx)
	// if err != nil {
	// 	fmt.Println("err1")
	// 	fmt.Println(err)
	// }
	msg, err := network.ReadMessage()
	if err != nil {
		fmt.Println("err2")
		fmt.Println(err)
	}
	fmt.Println(string(msg))

	err = network.SendAuthMessage()
	if err != nil {
		fmt.Println(err)
	}

	msg, err = network.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(msg))
	err = network.WriteMessage(NewLightOnRequest("group.living_room_lamps"))
	if err != nil {
		fmt.Println(err)
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
