package participant

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/agent"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

// Participant is a unified identity for any entity that can join a chat group.
// Exactly one of UserID or AgentID is set, unless Type is SYSTEM.
type Participant struct {
	ID          ID
	Type        ParticipantType
	UserID      *user.ID
	AgentID     *agent.ID
	DisplayName string
	AvatarURL   string
	CreatedAt   time.Time
}

func NewUserParticipant(userID user.ID, displayName, avatarURL string) *Participant {
	return &Participant{
		Type:        UserType,
		UserID:      &userID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}
}

func NewAgentParticipant(agentID agent.ID, displayName, avatarURL string) *Participant {
	return &Participant{
		Type:        AgentType,
		AgentID:     &agentID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}
}

func NewSystemParticipant() *Participant {
	return &Participant{
		Type:        SystemType,
		DisplayName: "System",
	}
}

func (p *Participant) IsUser() bool {
	return p.Type == UserType
}

func (p *Participant) IsAgent() bool {
	return p.Type == AgentType
}

func (p *Participant) IsSystem() bool {
	return p.Type == SystemType
}
