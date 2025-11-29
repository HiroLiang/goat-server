package repository

import (
	"context"
	"database/sql"

	"github.com/HiroLiang/goat-server/internal/database"
	"github.com/HiroLiang/goat-server/internal/database/model"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func (r *UserRepo) All(context *gin.Context) {

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
