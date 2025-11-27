package model

import "time"

type UserStatus string

const (
	UserStatusApplying  UserStatus = "APPLYING"
	UserStatusActivated UserStatus = "ACTIVATED"
	UserStatusDisabled  UserStatus = "DISABLED"
)

type User struct {
	ID         int64      `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"-"` // 不回傳密碼
	UserStatus UserStatus `db:"user_status" json:"user_status"`
	UserIP     string     `db:"user_ip" json:"user_ip"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

func (User) TableName() string {
	return "public.users"
}
