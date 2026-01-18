package websocket

import (
	"fmt"
	"log/slog"

	"saml.dev/gome-assistant/message"
)

// Subscription represents a websocket-level subscription to a
// particular message ID.
type Subscription struct {
	messageID int64
}

// MessageID returns the message ID that this subscription is
// subscribed to.
func (sub Subscription) MessageID() int64 {
	return sub.messageID
}

// Subscriber is called synchronously when a message is received that
// matches its subscription's message ID.
type Subscriber func(msg message.Message)

// NoopSubscriber is a `Subscriber` that does nothing.
func NoopSubscriber(_ message.Message) {}

// getSubscriber returns the subscriber, if any, that is subscribed to
// the specified message ID.
func (conn *Conn) getSubscriber(messageID int64) (Subscriber, bool) {
	conn.subscribersLock.RLock()
	defer conn.subscribersLock.RUnlock()

	subscriber, ok := conn.subscribers[messageID]
	return subscriber, ok
}

func (conn *Conn) SubscribeToEventType(eventType string, subr Subscriber) Subscription {
	var subn Subscription
	err := conn.Send(
		func(lc LockedConn) error {
			subn = lc.Subscribe(subr)
			e := message.SubscribeEventRequest{
				BaseMessage: message.BaseMessage{
					Type: "subscribe_events",
					ID:   subn.messageID,
				},
				EventType: eventType,
			}

			if err := lc.SendMessage(e); err != nil {
				lc.Unsubscribe(subn)
				return fmt.Errorf("error writing to websocket: %w", err)
			}
			// m, _ := ReadMessage(ctx, conn)
			// log.Default().Println(string(m))

			return nil
		},
	)

	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return subn
}

func (conn *Conn) SubscribeToStateChangedEvents(subr Subscriber) Subscription {
	return conn.SubscribeToEventType("state_changed", subr)
}
