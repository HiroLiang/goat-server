package chat

import (
	"encoding/json"
	"testing"

	"github.com/HiroLiang/goat-server/internal/interface/ws"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitTestEnv()
	m.Run()
}

func TestChatMessageHandler_Handle_ValidPayload(t *testing.T) {
	handler := NewMessageHandler()
	client := ws.NewClient(nil, nil, "user1")

	payload, _ := json.Marshal(SendPayload{RoomID: "room1", Content: "hello world"})
	err := handler.Handle(client, payload)

	assert.NoError(t, err)
}

func TestChatMessageHandler_Handle_EmptyContent(t *testing.T) {
	handler := NewMessageHandler()
	client := ws.NewClient(nil, nil, "user1")

	payload, _ := json.Marshal(SendPayload{RoomID: "room1", Content: ""})
	err := handler.Handle(client, payload)

	assert.NoError(t, err)
}

func TestChatMessageHandler_Handle_InvalidJSON_ReturnsError(t *testing.T) {
	handler := NewMessageHandler()
	client := ws.NewClient(nil, nil, "user1")

	err := handler.Handle(client, json.RawMessage(`not-valid-json`))

	assert.Error(t, err)
}

func TestChatMessageHandler_Handle_AnonymousClient(t *testing.T) {
	handler := NewMessageHandler()
	client := ws.NewClient(nil, nil, "") // no user ID

	payload, _ := json.Marshal(SendPayload{RoomID: "room1", Content: "hi"})
	err := handler.Handle(client, payload)

	assert.NoError(t, err)
}
