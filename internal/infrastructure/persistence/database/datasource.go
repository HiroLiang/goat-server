package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBName string

const (
	Postgres DBName = "postgres"
)

var dbNames = map[string]DBName{
	string(Postgres): Postgres,
}

var placeholders = map[DBName]squirrel.PlaceholderFormat{
	Postgres: squirrel.Dollar,
}

var dbs sync.Map

// InitDB -Initialize database map with config.
func InitDB() error {

	// Clear dbs
	dbs.Clear()

	// For each database config
	for name, conf := range config.App().Database {

		// check if db name is valid (set in enum DBName)
		if _, ok := isValidName(name); ok {

			// connect to database
			db, err := sqlx.Connect(conf.Driver, conf.Dsn)
			if err != nil {
				return fmt.Errorf("init database %s error: %v", name, err)
			}

			// set db config options
			if conf.Config != nil {
				db.SetMaxOpenConns(conf.Config.MaxOpenConns)
				db.SetMaxIdleConns(conf.Config.MaxIdleConns)
				db.SetConnMaxLifetime(time.Duration(conf.Config.ConnMaxLifetime) * time.Second)
				db.SetConnMaxIdleTime(time.Duration(conf.Config.ConnMaxIdleTime) * time.Second)
			}

			// try to ping the database
			if err := db.Ping(); err != nil {
				return fmt.Errorf("ping database %s error: %v", name, err)
			}

			// store database in the map
			dbs.Store(name, db)
			fmt.Printf("init database \"%s\" success\n", name)
		}
	}

	// Check if all databases are initialized
	if err := checkDBMap(); err != nil {
		return err
	}

	return nil
}

func GetDB(name DBName) *sqlx.DB {
	val, ok := dbs.Load(string(name))
	if !ok {
		return nil
	}
	return val.(*sqlx.DB)
}

func GetPlaceholder(name DBName) squirrel.PlaceholderFormat {
	return placeholders[name]
}

func CloseAllDBs() {
	dbs.Range(func(key, value interface{}) bool {
		if db, ok := value.(*sqlx.DB); ok {
			_ = db.Close()
			fmt.Println("Close DB connection: ", key.(string))
		}
		return true
	})
}

func isValidName(name string) (DBName, bool) {
	dbName, ok := dbNames[name]
	return dbName, ok
}

func checkDBMap() error {
	for name, key := range dbNames {
		if db := GetDB(key); db == nil {
			return fmt.Errorf("Database %s not initialized\n", name)
		}
	}

	return nil
}
