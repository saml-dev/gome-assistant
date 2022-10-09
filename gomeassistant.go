package gomeassistant

import (
	"context"
	"fmt"
	"time"

	"github.com/saml-dev/gome-assistant/internal/network"
	"nhooyr.io/websocket"
)

var ctx, ctxCancel = context.WithTimeout(context.Background(), time.Second*5)

var conn, _, err = websocket.Dial(ctx, "ws://192.168.86.67:8123/api/websocket", nil)

func main() {
	// sched := Schedule{
	// 	RunEvery: Daily,
	// }
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
