package repository

import (
	"github.com/HiroLiang/goat-server/internal/database"
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
