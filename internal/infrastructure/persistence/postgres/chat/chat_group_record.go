package chat

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type ChatGroupRecord struct {
	ID          chatgroup.ID        `db:"id"`
	Name        string              `db:"name"`
	Description string              `db:"description"`
	AvatarURL   string              `db:"avatar_url"`
	Type        chatgroup.GroupType `db:"type"`
	MaxMembers  int                 `db:"max_members"`
	IsDeleted   bool                `db:"is_deleted"`
	CreatedAt   time.Time           `db:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at"`
	CreatedBy   *user.ID            `db:"created_by"`
}
