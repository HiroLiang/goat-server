package ws

import "encoding/json"

// MessageHandler handles a specific type of WebSocket message.
type MessageHandler interface {
	Handle(client *Client, payload json.RawMessage) error
}
