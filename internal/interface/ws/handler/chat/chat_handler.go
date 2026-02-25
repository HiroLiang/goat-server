package chat

import (
	"encoding/json"

	"github.com/HiroLiang/goat-server/internal/interface/ws"
	"github.com/HiroLiang/goat-server/internal/logger"
	"go.uber.org/zap"
)

// SendPayload is the payload for a "chat.send" message.
type SendPayload struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

// MessageHandler handles "chat.send" messages.
type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handle(client *ws.Client, payload json.RawMessage) error {
	var p SendPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}

	logger.Log.Info("chat.send",
		zap.String("user_id", client.UserID),
		zap.String("room_id", p.RoomID),
		zap.String("content", p.Content),
	)

	// TODO: persist message and broadcast to room members via hub.SendToUser
	return nil
}
