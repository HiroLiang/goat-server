package chatmessage

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type ChatMessage struct {
	ID        ID
	GroupID   chatgroup.ID
	SenderID  participant.ID
	Content   string
	Type      MessageType
	ReplyToID *ID
	IsEdited  bool
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTextMessage(groupID chatgroup.ID, senderID participant.ID, content string) *ChatMessage {
	return &ChatMessage{
		GroupID:  groupID,
		SenderID: senderID,
		Content:  content,
		Type:     Text,
	}
}

func NewMessage(groupID chatgroup.ID, senderID participant.ID, content string, msgType MessageType) *ChatMessage {
	return &ChatMessage{
		GroupID:  groupID,
		SenderID: senderID,
		Content:  content,
		Type:     msgType,
	}
}

func (m *ChatMessage) Edit(content string) {
	m.Content = content
	m.IsEdited = true
}

func (m *ChatMessage) IsActive() bool {
	return !m.IsDeleted
}

func (m *ChatMessage) IsReply() bool {
	return m.ReplyToID != nil
}
