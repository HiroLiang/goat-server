package database

import "time"

type ConnectionConfig struct {
	Driver          string
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}
