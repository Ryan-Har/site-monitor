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
	AddMonitorResults(monitorResults ...MonitorResult) error
	GetMonitorResults(filters ...MonitorResultsFilter) ([]MonitorResult, error)
	DeleteMonitorResults(filters ...MonitorResultsFilter) error
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
	MonitorToSQLite(monitorTable string) (string, []interface{})
}

// ByMonitorIds implements Filter for filtering by monitorIds
type ByMonitorIds struct {
	Ids []int
}

// Apply implements Filter interface for ByMonitorIds
func (f ByMonitorIds) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Monitor_id IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ids))), placeholders
}

type ByUUIDs struct {
	Ids []string
}

func (f ByUUIDs) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.UUID IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ids))), placeholders
}

type ByUrls struct {
	Urls []string
}

func (f ByUrls) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Urls) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Urls))
	for i, id := range f.Urls {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Url IN (%s) ", monitorTable, generateQuestionMarks(len(f.Urls))), placeholders
}

type ByTypes struct {
	Types []string
}

func (f ByTypes) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Types) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Types))
	for i, id := range f.Types {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Type IN (%s) ", monitorTable, generateQuestionMarks(len(f.Types))), placeholders
}

type ByIntervalSecs struct {
	Intervals []int
}

func (f ByIntervalSecs) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Intervals) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Intervals))
	for i, id := range f.Intervals {
		placeholders[i] = id
	}
	fmt.Println(placeholders...)

	return fmt.Sprintf(" %s.Interval_in_seconds IN (%s) ", monitorTable, generateQuestionMarks(len(f.Intervals))), placeholders
}

type ByTimeoutSecs struct {
	Timeouts []int
}

func (f ByTimeoutSecs) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Timeouts) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Timeouts))
	for i, id := range f.Timeouts {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Timeout_in_seconds IN (%s) ", monitorTable, generateQuestionMarks(len(f.Timeouts))), placeholders
}

type ByPorts struct {
	Ports []int
}

func (f ByPorts) MonitorToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Ports) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Ports))
	for i, id := range f.Ports {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Port IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ports))), placeholders
}

func (h *SQLiteHandler) GetMonitors(filters ...MonitorFilter) ([]Monitor, error) {
	sqlTable := "Monitors"
	query := "SELECT Monitor_id, UUID, Url, Type, Interval_in_seconds, Timeout_in_seconds, Port FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.MonitorToSQLite(sqlTable)
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
		condSql, condArgs := filter.MonitorToSQLite(sqlTable)
		if condSql != "" {
			whereClause = append(whereClause, condSql)
			args = append(args, condArgs...)
		}
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	h.writeMutex.Lock()
	_, err := h.DB.Exec(query, args...)
	h.writeMutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}

type MonitorResult struct {
	CheckID        int
	MonitorID      int
	IsUp           int
	ResponseTimeMs int
	RunTimeEpoch   int
}

func (h *SQLiteHandler) AddMonitorResults(monitorResults ...MonitorResult) error {
	stmt, err := h.DB.Prepare("INSERT INTO Results (Monitor_id, Is_up, Response_time_in_ms, Run_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()
	for _, m := range monitorResults {
		_, err := stmt.Exec(m.MonitorID, m.IsUp, m.ResponseTimeMs, m.RunTimeEpoch)
		if err != nil {
			h.writeMutex.Unlock()
			return err
		}
	}
	h.writeMutex.Unlock()
	return nil
}

type MonitorResultsFilter interface {
	ResultsToSQLite(monitorTable string) (string, []interface{})
}

type ByCheckIds struct {
	Ids []int
}

func (f ByCheckIds) ResultsToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Check_id IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ids))), placeholders
}

func (f ByMonitorIds) ResultsToSQLite(monitorTable string) (string, []interface{}) {
	if len(f.Ids) == 0 {
		return "", nil
	}
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Monitor_id in (%s) ", monitorTable, generateQuestionMarks(len(f.Ids))), placeholders
}

type ByIsUp struct {
	Up bool
}

func (f ByIsUp) ResultsToSQLite(monitorTable string) (string, []interface{}) {

	placeholder := make([]interface{}, 1)
	if f.Up {
		placeholder[0] = 1
	} else {
		placeholder[0] = 0
	}

	return fmt.Sprintf(" %s.Is_up = ? ", monitorTable), placeholder
}

type BetweenRunTime struct {
	minEpoch int
	maxEpoch int
}

func (f BetweenRunTime) ResultsToSQLite(monitorTable string) (string, []interface{}) {

	placeholder := make([]interface{}, 2)
	placeholder[0] = f.minEpoch
	placeholder[0] = f.maxEpoch

	return fmt.Sprintf(" %s.Run_time BETWEEN ? AND ? ", monitorTable), placeholder
}

func (h *SQLiteHandler) GetMonitorResults(filters ...MonitorResultsFilter) ([]MonitorResult, error) {
	sqlTable := "Results"
	query := "SELECT Check_id, Monitor_id, Is_up, Response_time_in_ms, Run_time FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.ResultsToSQLite(sqlTable)
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

	var monitorResults []MonitorResult
	for rows.Next() {
		var m MonitorResult
		err := rows.Scan(&m.CheckID, &m.MonitorID, &m.IsUp, &m.ResponseTimeMs, &m.RunTimeEpoch)
		if err != nil {
			return nil, err
		}
		monitorResults = append(monitorResults, m)
	}
	return monitorResults, nil
}

func (h *SQLiteHandler) DeleteMonitorResults(filters ...MonitorResultsFilter) error {
	if len(filters) == 0 {
		return fmt.Errorf("filters must accompany delete monitors")
	}

	sqlTable := "Results"
	query := "DELETE FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.ResultsToSQLite(sqlTable)
		if condSql != "" {
			whereClause = append(whereClause, condSql)
			args = append(args, condArgs...)
		}
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	h.writeMutex.Lock()
	_, err := h.DB.Exec(query, args...)
	h.writeMutex.Unlock()
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
