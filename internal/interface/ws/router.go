package ws

import "fmt"

// MessageRouter dispatches incoming messages to the registered MessageHandler.
type MessageRouter struct {
	handlers map[string]MessageHandler
}

// NewMessageRouter creates a new MessageRouter.
func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers: make(map[string]MessageHandler),
	}
}

// Register binds msgType (e.g. "chat.send") to a handler.
func (r *MessageRouter) Register(msgType string, handler MessageHandler) {
	r.handlers[msgType] = handler
}

// Route dispatches msg to the appropriate handler.
// Returns an error if no handler is registered for msg.Type.
func (r *MessageRouter) Route(client *Client, msg *Message) error {
	handler, ok := r.handlers[msg.Type]
	if !ok {
		return fmt.Errorf("no handler for message type %q", msg.Type)
	}
	return handler.Handle(client, msg.Payload)
}
