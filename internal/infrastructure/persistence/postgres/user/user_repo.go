package user

import (
	"context"
	"errors"

	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var userColumns = []string{
	"id",
	"name",
	"email",
	"password",
	"user_status",
	"user_ip",
	"created_at",
	"updated_at",
}

var builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserRepository struct {
	db *sqlx.DB
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(dbName database.DBName) *UserRepository {
	return &UserRepository{db: database.GetDB(dbName)}
}

func (r *UserRepository) FindByID(ctx context.Context, id user.ID) (*user.User, error) {
	return r.findOneBy(ctx, squirrel.Eq{"id": id})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	return r.findOneBy(ctx, squirrel.Eq{"email": email})
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	rec := toRecord(u)

	query, args, err := builder.
		Insert("public.users").
		Columns("name", "email", "password", "user_status", "user_ip").
		Values(rec.Name, rec.Email, rec.Password, rec.UserStatus, rec.UserIP).
		ToSql()
	if err != nil {
		return err
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	rec := toRecord(u)

	query, args, err := builder.
		Update("public.users").
		Set("name", rec.Name).
		Set("password", rec.Password).
		Set("user_status", rec.UserStatus).
		Set("user_ip", rec.UserIP).
		Where(squirrel.Eq{"id": rec.ID}).
		ToSql()
	if err != nil {
		return err
	}

	return postgres.Exec(ctx, r.db, query, args...)
}

func baseSelect() squirrel.SelectBuilder {
	return builder.Select(userColumns...).From("public.users")
}

func (r *UserRepository) findOneBy(ctx context.Context, cond squirrel.Eq) (*user.User, error) {
	query, args, err := baseSelect().Where(cond).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	rec, err := postgres.ScanOne[UserRecord](ctx, r.db, query, args...)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(rec)
}
