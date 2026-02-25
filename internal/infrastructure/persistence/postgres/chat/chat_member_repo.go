package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/chatmember"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var ChatMemberTable = postgres.Table{
	Name: "public.chat_group_members",
	Columns: []string{
		"id",
		"group_id",
		"participant_id",
		"role",
		"joined_at",
		"is_archived",
		"is_muted",
		"is_pinned",
		"last_read_at",
		"updated_at",
	},
}

type ChatMemberRepository struct {
	db *sqlx.DB
}

var _ chatmember.Repository = (*ChatMemberRepository)(nil)

func NewChatMemberRepository(db *sqlx.DB) *ChatMemberRepository {
	return &ChatMemberRepository{db: db}
}

func (r *ChatMemberRepository) FindByID(ctx context.Context, id chatmember.ID) (*chatmember.ChatMember, error) {
	query, args, err := ChatMemberTable.Select(ChatMemberTable.Columns...).
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat member query: %w", err)
	}

	rec, err := postgres.ScanOne[ChatMemberRecord](ctx, r.db, query, args...)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, chatmember.ErrNotFound
		}
		return nil, fmt.Errorf("find chat member: %w", err)
	}

	return toChatMemberDomain(rec)
}

func (r *ChatMemberRepository) FindByGroupAndParticipant(
	ctx context.Context,
	groupID chatgroup.ID,
	participantID participant.ID,
) (*chatmember.ChatMember, error) {
	query, args, err := ChatMemberTable.Select(ChatMemberTable.Columns...).
		Where(squirrel.Eq{"group_id": groupID, "participant_id": participantID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat member query: %w", err)
	}

	rec, err := postgres.ScanOne[ChatMemberRecord](ctx, r.db, query, args...)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, chatmember.ErrNotFound
		}
		return nil, fmt.Errorf("find chat member: %w", err)
	}

	return toChatMemberDomain(rec)
}

func (r *ChatMemberRepository) FindByGroup(ctx context.Context, groupID chatgroup.ID) ([]*chatmember.ChatMember, error) {
	return r.findAll(ctx, squirrel.Eq{"group_id": groupID})
}

func (r *ChatMemberRepository) FindByParticipant(ctx context.Context, participantID participant.ID) ([]*chatmember.ChatMember, error) {
	return r.findAll(ctx, squirrel.Eq{"participant_id": participantID})
}

func (r *ChatMemberRepository) Add(ctx context.Context, m *chatmember.ChatMember) error {
	rec := toChatMemberRecord(m)

	query, args, err := ChatMemberTable.Insert().
		Columns("group_id", "participant_id", "role").
		Values(rec.GroupID, rec.ParticipantID, rec.Role).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert chat member: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *ChatMemberRepository) Update(ctx context.Context, m *chatmember.ChatMember) error {
	rec := toChatMemberRecord(m)

	query, args, err := ChatMemberTable.Update().
		Set("role", rec.Role).
		Set("is_archived", rec.IsArchived).
		Set("is_muted", rec.IsMuted).
		Set("is_pinned", rec.IsPinned).
		Set("last_read_at", rec.LastReadAt).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": rec.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update chat member: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *ChatMemberRepository) Remove(ctx context.Context, groupID chatgroup.ID, participantID participant.ID) error {
	query, args, err := ChatMemberTable.Delete().
		Where(squirrel.Eq{"group_id": groupID, "participant_id": participantID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build remove chat member: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *ChatMemberRepository) findAll(
	ctx context.Context,
	cond squirrel.Sqlizer,
) ([]*chatmember.ChatMember, error) {

	query, args, err := ChatMemberTable.Select(ChatMemberTable.Columns...).
		Where(cond).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat members query: %w", err)
	}

	records, err := postgres.ScanAll[ChatMemberRecord](ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan chat members: %w", err)
	}

	members := make([]*chatmember.ChatMember, 0, len(records))
	for _, rec := range records {
		m, err := toChatMemberDomain(&rec)
		if err != nil {
			return nil, fmt.Errorf("convert chat member: %w", err)
		}
		members = append(members, m)
	}

	return members, nil
}
