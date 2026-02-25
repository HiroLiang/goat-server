package ws

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageRouter_Route_DispatchesToHandler(t *testing.T) {
	router := NewMessageRouter()
	handler := &mockHandler{}
	router.Register("chat.send", handler)

	payload := json.RawMessage(`{"room_id":"r1","content":"hi"}`)
	err := router.Route(nil, &Message{Type: "chat.send", Payload: payload})

	assert.NoError(t, err)
	assert.Equal(t, 1, handler.Calls())
	assert.Equal(t, payload, handler.payload)
}

func TestMessageRouter_Route_UnknownType_ReturnsError(t *testing.T) {
	router := NewMessageRouter()

	err := router.Route(nil, &Message{Type: "no.such.type"})

	assert.ErrorContains(t, err, "no.such.type")
}

func TestMessageRouter_Route_PropagatesHandlerError(t *testing.T) {
	router := NewMessageRouter()
	want := errors.New("handler failed")
	router.Register("game.move", &mockHandler{err: want})

	err := router.Route(nil, &Message{Type: "game.move"})

	assert.ErrorIs(t, err, want)
}

func TestMessageRouter_Register_LastHandlerWins(t *testing.T) {
	router := NewMessageRouter()
	first := &mockHandler{}
	second := &mockHandler{}
	router.Register("chat.send", first)
	router.Register("chat.send", second) // overwrites

	_ = router.Route(nil, &Message{Type: "chat.send"})

	assert.Equal(t, 0, first.Calls(), "first handler should be replaced")
	assert.Equal(t, 1, second.Calls())
}

func TestMessageRouter_Route_PassesClientToHandler(t *testing.T) {
	router := NewMessageRouter()
	var gotClient *Client
	router.Register("ping", &captureClientHandler{fn: func(c *Client) { gotClient = c }})

	client := newTestClient(nil, "user1")
	_ = router.Route(client, &Message{Type: "ping"})

	assert.Same(t, client, gotClient)
}

// captureClientHandler is a one-off helper for testing client propagation.
type captureClientHandler struct {
	fn func(*Client)
}

func (h *captureClientHandler) Handle(c *Client, _ json.RawMessage) error {
	h.fn(c)
	return nil
}
