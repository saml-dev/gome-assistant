package websocket

// Subscriber is called synchronously when a message is received that
// matches its subscription.
type Subscriber func(msg ChanMsg)
