package database

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/HiroLiang/goat-chat-server/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

type DBName string

const (
	SQLite DBName = "sqlite"
)

var dbNames = map[string]DBName{
	string(SQLite): SQLite,
}

var dbs sync.Map

func InitDB() error {
	for name, conf := range config.Cfg.Database {
		if _, ok := isValidName(name); ok {
			conn, err := sql.Open(conf.Driver, conf.Dsn)

			if err != nil {
				return fmt.Errorf("init database %s error: %v", name, err)
			}

			if err := conn.Ping(); err != nil {
				return fmt.Errorf("ping database %s error: %v", name, err)
			}

			dbs.Store(name, conn)
			fmt.Printf("init database \"%s\" success\n", name)
		} else {
			return fmt.Errorf("invalid db name: %s", name)
		}
	}
	return nil
}

func GetDB(name DBName) (*sql.DB, bool) {
	val, ok := dbs.Load(string(name))
	if !ok {
		return nil, false
	}
	return val.(*sql.DB), true
}

func CloseAllDBs() {
	dbs.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*sql.DB); ok {
			_ = conn.Close()
			fmt.Println("Close DB connection: ", key.(string))
		}
		return true
	})
}

func isValidName(name string) (DBName, bool) {
	dbName, ok := dbNames[name]
	return dbName, ok
}
