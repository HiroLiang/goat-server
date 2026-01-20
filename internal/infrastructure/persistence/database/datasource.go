package database

import (
	"fmt"
	"sync"
	"time"

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

type DataSources struct {
	Sources sync.Map
}

func NewDataSources(databases map[string]ConnectionConfig) (*DataSources, error) {

	// Clear dbs
	var dataSources = &DataSources{}

	// For each database config
	for name, conf := range databases {

		// check if db name is valid (set in enum DBName)
		if _, ok := isValidName(name); ok {

			// connect to database
			db, err := sqlx.Connect(conf.Driver, conf.Dsn)
			if err != nil {
				return nil, fmt.Errorf("init database %s error: %v", name, err)
			}

			// set db config options
			if conf.MaxOpenConns > 0 {
				db.SetMaxOpenConns(conf.MaxOpenConns)
			}
			if conf.MaxIdleConns > 0 {
				db.SetMaxIdleConns(conf.MaxIdleConns)
			}
			if conf.ConnMaxLifetime > 0 {
				db.SetConnMaxLifetime(conf.ConnMaxLifetime * time.Second)
			}
			if conf.ConnMaxIdleTime > 0 {
				db.SetConnMaxIdleTime(conf.ConnMaxIdleTime * time.Second)
			}

			// try to ping the database
			if err := db.Ping(); err != nil {
				return nil, fmt.Errorf("ping database %s error: %v", name, err)
			}

			// store database in the map
			dataSources.Sources.Store(name, db)
			fmt.Printf("init database \"%s\" success\n", name)
		}
	}

	// Check if all databases are initialized
	if err := checkDBMap(dataSources); err != nil {
		return nil, err
	}

	return dataSources, nil
}

func (sources *DataSources) GetDB(name DBName) *sqlx.DB {
	val, ok := sources.Sources.Load(string(name))
	if !ok {
		return nil
	}
	return val.(*sqlx.DB)
}

func (sources *DataSources) GetPlaceholder(name DBName) squirrel.PlaceholderFormat {
	return placeholders[name]
}

func (sources *DataSources) CloseAllDBs() {
	sources.Sources.Range(func(key, value interface{}) bool {
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

func checkDBMap(dataSources *DataSources) error {
	for name, key := range dbNames {
		if db := dataSources.GetDB(key); db == nil {
			return fmt.Errorf("Database %s not initialized\n", name)
		}
	}

	return nil
}
