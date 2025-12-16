package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ScanOne[T any](ctx context.Context, db *sqlx.DB, query string, args ...any) (*T, error) {
	var rec T
	if err := db.GetContext(ctx, &rec, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan one: %w", err)
	}
	return &rec, nil
}

func ScanAll[T any](ctx context.Context, db *sqlx.DB, query string, args ...any) ([]T, error) {
	var list []T
	if err := db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, fmt.Errorf("scan all: %w", err)
	}
	return list, nil
}

func Exec(ctx context.Context, db *sqlx.DB, query string, args ...any) error {
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	return nil
}
