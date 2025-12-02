package user

import "time"

type StatusDB string

const (
	Active   StatusDB = "active"
	Inactive StatusDB = "inactive"
	Banned   StatusDB = "banned"
	Applying StatusDB = "applying"
	Deleted  StatusDB = "deleted"
)

type UserRecord struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name" `
	Email      string    `db:"email" `
	Password   string    `db:"password"`
	UserStatus StatusDB  `db:"user_status"`
	UserIP     string    `db:"user_ip"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
