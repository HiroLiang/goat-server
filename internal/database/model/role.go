package model

import "time"

type Role struct {
	ID        uint      `gorm:"primary_key;column:id"`
	Name      string    `gorm:"column:name;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Role) TableName() string {
	return "public.roles"
}
