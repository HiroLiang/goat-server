package role

import "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres"

var Table = postgres.Table{
	Name: "goat.public.roles",
	Columns: []string{
		"id",
		"type",
		"creator",
		"created_at",
		"updated_at",
	},
}
