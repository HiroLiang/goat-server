package chatmember

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type ChatMember struct {
	ID            ID
	GroupID       chatgroup.ID
	ParticipantID participant.ID
	Role          Role
	JoinedAt      time.Time
	IsArchived    bool
	IsMuted       bool
	IsPinned      bool
	LastReadAt    *time.Time
	UpdatedAt     time.Time
}

func NewChatMember(groupID chatgroup.ID, participantID participant.ID, role Role) *ChatMember {
	return &ChatMember{
		GroupID:       groupID,
		ParticipantID: participantID,
		Role:          role,
	}
}

func (m *ChatMember) IsOwner() bool {
	return m.Role == Owner
}

func (m *ChatMember) IsAdmin() bool {
	return m.Role == Admin
}

func (m *ChatMember) CanManageMembers() bool {
	return m.Role == Owner || m.Role == Admin
}

func (m *ChatMember) MarkAsRead(at time.Time) {
	m.LastReadAt = &at
}
