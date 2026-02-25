package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/chatmessage"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var CharMessageTable = postgres.Table{
	Name: "public.chat_records",
	Columns: []string{
		"id",
		"group_id",
		"sender_id",
		"content",
		"message_type",
		"reply_to_id",
		"is_edited",
		"is_deleted",
		"created_at",
		"updated_at",
	},
}

type ChatMessageRepository struct {
	db *sqlx.DB
}

var _ chatmessage.Repository = (*ChatMessageRepository)(nil)

func NewChatMessageRepository(db *sqlx.DB) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

func (r *ChatMessageRepository) FindByID(ctx context.Context, id chatmessage.ID) (*chatmessage.ChatMessage, error) {
	query, args, err := CharMessageTable.Select(CharMessageTable.Columns...).
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat message query: %w", err)
	}

	rec, err := postgres.ScanOne[ChatMessageRecord](ctx, r.db, query, args...)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, chatmessage.ErrNotFound
		}
		return nil, fmt.Errorf("find chat message: %w", err)
	}

	return toChatMessageDomain(rec)
}

func (r *ChatMessageRepository) FindByGroup(
	ctx context.Context,
	groupID chatgroup.ID,
	limit, offset uint64,
) ([]*chatmessage.ChatMessage, error) {
	query, args, err := CharMessageTable.Select(CharMessageTable.Columns...).
		Where(squirrel.Eq{"group_id": groupID}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat messages query: %w", err)
	}

	records, err := postgres.ScanAll[ChatMessageRecord](ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan chat messages: %w", err)
	}

	messages := make([]*chatmessage.ChatMessage, 0, len(records))
	for _, rec := range records {
		msg, err := toChatMessageDomain(&rec)
		if err != nil {
			return nil, fmt.Errorf("convert chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *ChatMessageRepository) FindByGroupBefore(
	ctx context.Context,
	groupID chatgroup.ID,
	beforeID chatmessage.ID,
	limit uint64,
) ([]*chatmessage.ChatMessage, error) {
	query, args, err := CharMessageTable.Select(CharMessageTable.Columns...).
		Where(squirrel.And{
			squirrel.Eq{"group_id": groupID},
			squirrel.Eq{"is_deleted": false},
			squirrel.Lt{"id": beforeID},
		}).
		OrderBy("id DESC").
		Limit(limit).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat messages before query: %w", err)
	}

	records, err := postgres.ScanAll[ChatMessageRecord](ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan chat messages before: %w", err)
	}

	messages := make([]*chatmessage.ChatMessage, 0, len(records))
	for _, rec := range records {
		msg, err := toChatMessageDomain(&rec)
		if err != nil {
			return nil, fmt.Errorf("convert chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	// Reverse to ascending order (oldest first in the batch)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatMessageRepository) FindLatestByGroup(
	ctx context.Context,
	groupID chatgroup.ID,
) (*chatmessage.ChatMessage, error) {
	query, args, err := CharMessageTable.Select(CharMessageTable.Columns...).
		Where(squirrel.Eq{"group_id": groupID, "is_deleted": false}).
		OrderBy("id DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build latest message query: %w", err)
	}

	rec, err := postgres.ScanOne[ChatMessageRecord](ctx, r.db, query, args...)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, chatmessage.ErrNotFound
		}
		return nil, fmt.Errorf("find latest message: %w", err)
	}

	return toChatMessageDomain(rec)
}

func (r *ChatMessageRepository) CountByGroupAfter(
	ctx context.Context,
	groupID chatgroup.ID,
	since time.Time,
) (int64, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("COUNT(*)").
		From(CharMessageTable.Name).
		Where(squirrel.And{
			squirrel.Eq{"group_id": groupID, "is_deleted": false},
			squirrel.Gt{"created_at": since},
		}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count messages after: %w", err)
	}

	return count, nil
}

func (r *ChatMessageRepository) FindBySender(
	ctx context.Context,
	senderID participant.ID,
) ([]*chatmessage.ChatMessage, error) {
	query, args, err := CharMessageTable.Select(CharMessageTable.Columns...).
		Where(squirrel.Eq{"sender_id": senderID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build chat messages query: %w", err)
	}

	records, err := postgres.ScanAll[ChatMessageRecord](ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan chat messages: %w", err)
	}

	messages := make([]*chatmessage.ChatMessage, 0, len(records))
	for _, rec := range records {
		msg, err := toChatMessageDomain(&rec)
		if err != nil {
			return nil, fmt.Errorf("convert chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *ChatMessageRepository) Create(ctx context.Context, msg *chatmessage.ChatMessage) error {
	rec := toChatMessageRecord(msg)

	query, args, err := CharMessageTable.Insert().
		Columns("group_id", "sender_id", "content", "message_type", "reply_to_id").
		Values(rec.GroupID, rec.SenderID, rec.Content, rec.Type, rec.ReplyToID).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert chat message: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *ChatMessageRepository) Update(ctx context.Context, msg *chatmessage.ChatMessage) error {
	rec := toChatMessageRecord(msg)

	query, args, err := CharMessageTable.Update().
		Set("content", rec.Content).
		Set("is_edited", rec.IsEdited).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": rec.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update chat message: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *ChatMessageRepository) SoftDelete(ctx context.Context, id chatmessage.ID) error {
	query, args, err := CharMessageTable.Update().
		Set("is_deleted", true).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build soft delete chat message: %w", err)
	}

	return postgres.Exec(ctx, r.db, query, args...)
}
