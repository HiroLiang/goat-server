package model

type Role struct {
	ID          uint   `gorm:"primary_key;column:id"`
	Name        string `gorm:"column:name;not null"`
	Permissions string `gorm:"column:permissions;not null"`
}

func (Role) TableName() string {
	return "public.roles"
}
