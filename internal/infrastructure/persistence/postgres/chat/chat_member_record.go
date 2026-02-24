package chat

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/chatmember"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type ChatMemberRecord struct {
	ID            chatmember.ID   `db:"id"`
	GroupID       chatgroup.ID    `db:"group_id"`
	ParticipantID participant.ID  `db:"participant_id"`
	Role          chatmember.Role `db:"role"`
	JoinedAt      time.Time       `db:"joined_at"`
	IsArchived    bool            `db:"is_archived"`
	IsMuted       bool            `db:"is_muted"`
	IsPinned      bool            `db:"is_pinned"`
	LastReadAt    *time.Time      `db:"last_read_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}
