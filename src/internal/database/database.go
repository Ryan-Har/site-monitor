package database

import (
	"database/sql"
	"fmt"

	_ "embed"
	"log/slog"
	"strings"
	"sync"

	"github.com/Ryan-Har/site-monitor/src/config"
	_ "github.com/mattn/go-sqlite3"
)

type DBHandler interface {
	AddMonitors(monitor ...Monitor) error
	GetMonitors(filters ...MonitorFilter) ([]Monitor, error)
	DeleteMonitors(filters ...MonitorFilter) error
	AddMonitorResults(monitorResults ...MonitorResult) error
	GetMonitorResults(filters ...MonitorResultsFilter) ([]MonitorResult, error)
	DeleteMonitorResults(filters ...MonitorResultsFilter) error
	AddNotification(notificationSettings NotificationSettings) error
	GetNotifications(filters ...NotificationFilter) ([]NotificationSettings, error)
	DeleteNotifications(filters ...NotificationFilter) error
}

type SQLiteHandler struct {
	DB         *sql.DB
	Version    int32 //version of database schema
	writeMutex *sync.Mutex
}

func NewSQLiteHandler() (*DBHandler, error) {
	var dbHandler DBHandler
	dbLoc := config.GetConfig().SQLITE_DB_LOCATION
	db, err := openSQLiteDB(dbLoc)
	if err != nil {
		return &dbHandler, err
	}
	if !isSQLiteDBPopulated(db) {
		slog.Info("db in new state, applying schema")
		if err = populateSQLiteDB(db); err != nil {
			return &dbHandler, fmt.Errorf("unable to apply db schema. %s", err.Error())
		}

	}
	dbHandler = &SQLiteHandler{
		DB:         db,
		Version:    1,
		writeMutex: &sync.Mutex{},
	}
	return &dbHandler, nil
}

func openSQLiteDB(dbLoc string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		// Handle error setting foreign_keys
		return nil, fmt.Errorf("unable to enforce foreign keys in db: %v", err.Error())
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

	return rows.Next()
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

func generateQuestionMarks(n int) string {
	if n <= 0 {
		return ""
	}

	questionMarks := make([]string, n)
	for i := 0; i < n; i++ {
		questionMarks[i] = "?"
	}

	return strings.Join(questionMarks, ",")
}
