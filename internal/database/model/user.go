package model

import "time"

type User struct {
	ID        uint      `gorm:"primary_key;column:id"`
	Name      string    `gorm:"column:name;not null"`
	Account   string    `gorm:"column:account;not null"`
	Password  string    `gorm:"column:password;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "public.users"
}
