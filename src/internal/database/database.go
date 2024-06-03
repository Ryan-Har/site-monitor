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
	AddMonitors(monitor ...Monitor) error
	GetMonitors(filters ...MonitorFilter) ([]Monitor, error)
	DeleteMonitors(filters ...MonitorFilter) error
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

func (h *SQLiteHandler) AddMonitors(monitor ...Monitor) error {
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

// Filter defines the interface for filter functions
type MonitorFilter interface {
	ToSQLite(monitorTable string) (string, []interface{})
}

// ByMonitorIds implements Filter for filtering by monitorIds
type ByMonitorIds struct {
	ids []int
}

// Apply implements Filter interface for ByMonitorIds
func (f ByMonitorIds) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.ids))
	for i, id := range f.ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Monitor_id IN (%s) ", monitorTable, strings.Repeat("?", len(f.ids))), placeholders
}

type ByUUIDs struct {
	ids []string
}

func (f ByUUIDs) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.ids))
	for i, id := range f.ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.UUID IN (%s) ", monitorTable, strings.Repeat("?", len(f.ids))), placeholders
}

type ByUrls struct {
	urls []string
}

func (f ByUrls) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.urls) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.urls))
	for i, id := range f.urls {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Url IN (%s) ", monitorTable, strings.Repeat("?", len(f.urls))), placeholders
}

type ByTypes struct {
	types []string
}

func (f ByTypes) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.types) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.types))
	for i, id := range f.types {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Type IN (%s) ", monitorTable, strings.Repeat("?", len(f.types))), placeholders
}

type ByIntervalSecs struct {
	intervals []int
}

func (f ByIntervalSecs) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.intervals) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.intervals))
	for i, id := range f.intervals {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Interval_in_seconds IN (%s) ", monitorTable, strings.Repeat("?", len(f.intervals))), placeholders
}

type ByTimeoutSecs struct {
	timeouts []int
}

func (f ByTimeoutSecs) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.timeouts) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.timeouts))
	for i, id := range f.timeouts {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Timeout_in_seconds IN (%s) ", monitorTable, strings.Repeat("?", len(f.timeouts))), placeholders
}

type ByPorts struct {
	ports []int
}

func (f ByPorts) ToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.ports) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.ports))
	for i, id := range f.ports {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Port IN (%s) ", monitorTable, strings.Repeat("?", len(f.ports))), placeholders
}

func (h *SQLiteHandler) GetMonitors(filters ...MonitorFilter) ([]Monitor, error) {
	sqlTable := "Monitors"
	query := "SELECT Monitor_id, UUID, Url, Type, Interval_in_seconds, Timeout_in_seconds, Port FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.ToSQLite(sqlTable)
		if condSql != "" {
			whereClause = append(whereClause, condSql)
			args = append(args, condArgs...)
		}
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []Monitor
	for rows.Next() {
		var m Monitor
		err := rows.Scan(&m.MonitorID, &m.UUID, &m.URL, &m.Type, &m.IntervalSecs, &m.TimeoutSecs, &m.Port)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	return monitors, nil
}

func (h *SQLiteHandler) DeleteMonitors(filters ...MonitorFilter) error {
	if len(filters) == 0 {
		return fmt.Errorf("filters must accompany delete monitors")
	}

	sqlTable := "Monitors"
	query := "DELETE FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.ToSQLite(sqlTable)
		if condSql != "" {
			whereClause = append(whereClause, condSql)
			args = append(args, condArgs...)
		}
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	_, err := h.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
