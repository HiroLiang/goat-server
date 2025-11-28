package model

import "time"

/*
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
*/

type User struct {
	ID         int64     `db:"id" gorm:"primary_key;column:id"`
	Name       string    `db:"name" gorm:"column:name;not null"`
	Email      string    `db:"email" gorm:"column:email;not null"`
	Password   string    `db:"password" gorm:"column:password;not null"`
	UserStatus string    `db:"user_status" gorm:"column:user_status;not null"`
	UserIP     string    `db:"user_ip" gorm:"column:user_ip;not null"`
	CreatedAt  time.Time `db:"created_at" gorm:"column:created_at"`
}

func (User) TableName() string {
	return "public.users"
}
