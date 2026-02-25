package chat

import "github.com/HiroLiang/goat-server/internal/domain/chatmember"

func toChatMemberDomain(rec *ChatMemberRecord) (*chatmember.ChatMember, error) {
	return &chatmember.ChatMember{
		ID:            rec.ID,
		GroupID:       rec.GroupID,
		ParticipantID: rec.ParticipantID,
		Role:          rec.Role,
		JoinedAt:      rec.JoinedAt,
		IsArchived:    rec.IsArchived,
		IsMuted:       rec.IsMuted,
		IsPinned:      rec.IsPinned,
		LastReadAt:    rec.LastReadAt,
		UpdatedAt:     rec.UpdatedAt,
	}, nil
}

func toChatMemberRecord(m *chatmember.ChatMember) *ChatMemberRecord {
	return &ChatMemberRecord{
		ID:            m.ID,
		GroupID:       m.GroupID,
		ParticipantID: m.ParticipantID,
		Role:          m.Role,
		JoinedAt:      m.JoinedAt,
		IsArchived:    m.IsArchived,
		IsMuted:       m.IsMuted,
		IsPinned:      m.IsPinned,
		LastReadAt:    m.LastReadAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
