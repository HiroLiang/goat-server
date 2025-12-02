package role

import "time"

type Role struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	CreateBy  string    `db:"create_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
