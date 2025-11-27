package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/database"
	"github.com/HiroLiang/goat-server/internal/database/model"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDisabled      = errors.New("user is disabled")
)

type UserRepo struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepo(dbName database.DBName) *UserRepo {
	return &UserRepo{
		db:      database.GetDB(dbName),
		builder: squirrel.StatementBuilder.PlaceholderFormat(database.GetPlaceholder(dbName)),
	}
}

// GetByEmail 根據 email 查詢用戶
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := r.builder.
		Select("id", "name", "email", "password", "user_status", "user_ip", "created_at").
		From("public.users").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query error: %w", err)
	}

	var user model.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user error: %w", err)
	}

	return &user, nil
}

// GetByID 根據 ID 查詢用戶
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query, args, err := r.builder.
		Select("id", "name", "email", "password", "user_status", "user_ip", "created_at").
		From("public.users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query error: %w", err)
	}

	var user model.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user error: %w", err)
	}

	return &user, nil
}

// Create 創建新用戶
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query, args, err := r.builder.
		Insert("public.users").
		Columns("name", "email", "password", "user_status", "user_ip").
		Values(user.Name, user.Email, user.Password, user.UserStatus, user.UserIP).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("build query error: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}

// UpdateUserIP 更新用戶 IP
func (r *UserRepo) UpdateUserIP(ctx context.Context, userID int64, ip string) error {
	query, args, err := r.builder.
		Update("public.users").
		Set("user_ip", ip).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update user IP error: %w", err)
	}

	return nil
}

// UpdateStatus 更新用戶狀態
func (r *UserRepo) UpdateStatus(ctx context.Context, userID int64, status model.UserStatus) error {
	query, args, err := r.builder.
		Update("public.users").
		Set("user_status", status).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update user status error: %w", err)
	}

	return nil
}

// IsActivated 檢查用戶是否已激活
func (r *UserRepo) IsActivated(user *model.User) bool {
	return user.UserStatus == model.UserStatusActivated
}

// IsDisabled 檢查用戶是否已停用
func (r *UserRepo) IsDisabled(user *model.User) bool {
	return user.UserStatus == model.UserStatusDisabled
}

// EmailExists 檢查 email 是否已存在
func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	query, args, err := r.builder.
		Select("1").
		From("public.users").
		Where(squirrel.Eq{"email": email}).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("build query error: %w", err)
	}

	var exists int
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("query email exists error: %w", err)
	}

	return true, nil
}

func (r *UserRepo) ExistsApplyingByIP(ctx context.Context, ip string) (bool, error) {
	query, args, err := r.builder.
		Select("1").
		From("public.users").
		Where(squirrel.Eq{
			"user_ip":     ip,
			"user_status": "APPLYING",
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var dummy int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&dummy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepo) CreateApplyingUser(ctx context.Context, email, hashedPassword, ip string) (*model.User, error) {
	query, args, err := r.builder.
		Insert("public.users").
		Columns("name", "email", "password", "user_status", "user_ip").
		Values("", email, hashedPassword, "APPLYING", ip).
		Suffix("RETURNING id, name, email, password, user_status, user_ip, created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = r.db.QueryRowxContext(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.UserStatus,
		&u.UserIP,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) UpdateUserName(ctx context.Context, id int64, name string) error {
	query, args, err := r.builder.
		Update("public.users").
		Set("name", name).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
