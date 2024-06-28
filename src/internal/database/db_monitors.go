package database

import (
	"fmt"
	"strings"
)

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
	defer h.writeMutex.Unlock()
	for _, m := range monitor {
		_, err := stmt.Exec(m.UUID, m.URL, m.Type, m.IntervalSecs, m.TimeoutSecs, m.Port)
		if err != nil {
			return err
		}
	}
	return nil
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
