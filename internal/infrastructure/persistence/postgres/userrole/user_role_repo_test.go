package userrole

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/testutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// TestUserRoleRepository_Exists Test is the Exist query correct
func TestUserRoleRepository_Exists(t *testing.T) {
	db, mock := testutil.SetupDB(t)
	repo := UserRoleRepository{db: sqlx.NewDb(db, "postgres")}

	mock.ExpectQuery(`SELECT 1 .*users_roles.*roles.*`).
		WithArgs("u1", "admin").
		WillReturnRows(
			sqlmock.NewRows([]string{"1"}).AddRow(1),
		)

	ok := repo.Exists(context.Background(), "u1", role.Admin)
	assert.True(t, ok)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRoleRepository_Assign Test is the Assign query correct
func TestUserRoleRepository_Assign(t *testing.T) {
	db, mock := testutil.SetupDB(t)
	repo := UserRoleRepository{db: sqlx.NewDb(db, "postgres")}

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO goat.public.users_roles (user_id,role_id)
		SELECT $1, id FROM goat.public.roles WHERE type = $2
		ON CONFLICT DO NOTHING
	`)).
		WithArgs("u1", "admin").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Assign(context.Background(), "u1", role.Admin)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRoleRepository_Revoke(t *testing.T) {
	db, mock := testutil.SetupDB(t)
	repo := UserRoleRepository{db: sqlx.NewDb(db, "postgres")}

	//noinspection SqlWithoutWhere
	mock.ExpectExec(
		`DELETE FROM.*users_roles.* WHERE .* USING .*roles.*`,
	).
		WithArgs("admin", "u1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Revoke(context.Background(), "u1", role.Admin)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
