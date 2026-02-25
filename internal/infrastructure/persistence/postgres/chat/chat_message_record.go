package chat

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/chatmessage"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type ChatMessageRecord struct {
	ID        chatmessage.ID          `db:"id"`
	GroupID   chatgroup.ID            `db:"group_id"`
	SenderID  participant.ID          `db:"sender_id"`
	Content   string                  `db:"content"`
	Type      chatmessage.MessageType `db:"message_type"`
	ReplyToID *chatmessage.ID         `db:"reply_to_id"`
	IsEdited  bool                    `db:"is_edited"`
	IsDeleted bool                    `db:"is_deleted"`
	CreatedAt time.Time               `db:"created_at"`
	UpdatedAt time.Time               `db:"updated_at"`
}
