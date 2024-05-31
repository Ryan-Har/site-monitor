package db

import (
	"database/sql"
	"fmt"

	_ "embed"
	"github.com/Ryan-Har/site-monitor/src/config"
	_ "github.com/mattn/go-sqlite3"
)

type DBHandler interface {
}

type SQLiteHandler struct {
	DB      *sql.DB
	Version int32 //version of database schema
}

func NewSQLiteHandler() (*DBHandler, error) {
	var dbHandler DBHandler
	db, err := openSQLiteDB()
	if err != nil {
		return &dbHandler, err
	}
	if !isSQLiteDBPopulated(db) {
		fmt.Println("db in new state, applying schema")
		if err = populateSQLiteDB(db); err != nil {
			return &dbHandler, fmt.Errorf("unable to apply db schema. %s", err.Error())
		}

	}
	dbHandler = SQLiteHandler{
		DB:      db,
		Version: 1,
	}
	return &dbHandler, nil
}

func openSQLiteDB() (*sql.DB, error) {
	dbLoc := config.GetConfig().SQLITE_DB_LOCATION

	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func isSQLiteDBPopulated(db *sql.DB) bool {
	tableName := "Monitors" //just checking if this table exists since it's so crucial

	rows, err := db.Query("SELECT name FROM PRAGMA_TABLE_INFO (?)", tableName)
	if err != nil {
		return false
	}

	defer rows.Close()

	if rows.Next() {
		fmt.Println("row exists")
		return true
	}
	return false
}

//go:embed db_schema.sql
var sqliteSchema []byte

func populateSQLiteDB(db *sql.DB) error {
	schemaFile := string(sqliteSchema)

	_, err := db.Exec(schemaFile)
	if err != nil {
		return err
	}
	return nil
}
