package database

import (
	"database/sql"
	"fmt"

	_ "embed"
	"github.com/Ryan-Har/site-monitor/src/config"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
)

type DBHandler interface {
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
		fmt.Println("db in new state, applying schema")
		if err = populateSQLiteDB(db); err != nil {
			return &dbHandler, fmt.Errorf("unable to apply db schema. %s", err.Error())
		}

	}
	dbHandler = SQLiteHandler{
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

type Monitor struct {
	MonitorID    int
	UUID         string
	URL          string
	Type         string
	IntervalSecs int
	TimeoutSecs  int
	Port         int
}

func (h *SQLiteHandler) AddMonitor(monitor ...Monitor) error {
	stmt, err := h.DB.Prepare("INSERT INTO Monitors (UUID, Url, Type, Interval_in_seconds, Timeout_In_Seconds, Port) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()
	for _, m := range monitor {
		_, err := stmt.Exec(m.UUID, m.URL, m.Type, m.IntervalSecs, m.TimeoutSecs, m.Port)
		if err != nil {
			h.writeMutex.Unlock()
			return err
		}
	}
	h.writeMutex.Unlock()
	return nil
}

func (h *SQLiteHandler) GetMonitorByID(ids ...int) ([]Monitor, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("no id provided, nothing to get")
	}

	stmt, err := h.DB.Prepare("SELECT Monitor_id, UUID, Url, Type, Interval_in_seconds, Timeout_in_seconds, Port FROM Monitors WHERE Monitor_id in (?" + strings.Repeat(",?", len(ids)-1) + ")")
	if err != nil {
		return nil, err
	}

	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}

	rows, err := stmt.Query(idsInterface...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitorsResult := []Monitor{}

	for rows.Next() {
		var m Monitor
		if err := rows.Scan(&m.MonitorID, &m.UUID, &m.URL, &m.Type, &m.IntervalSecs, &m.TimeoutSecs, &m.Port); err != nil {
			return nil, err
		}
		monitorsResult = append(monitorsResult, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return monitorsResult, nil
}

func (h *SQLiteHandler) DeleteMonitorByID(ids ...int) error {
	if len(ids) == 0 {
		return fmt.Errorf("no id provided, nothing to delete")
	}

	stms, err := h.DB.Prepare("DELETE FROM Monitors WHERE Monitor_id in (?" + strings.Repeat(",?", len(ids)-1) + ")")
	if err != nil {
		return err
	}

	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}

	_, err = stms.Exec(idsInterface...)
	if err != nil {
		return err
	}

	return nil
}
