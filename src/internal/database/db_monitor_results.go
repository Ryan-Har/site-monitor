package database

import (
	"fmt"
	"strings"
)

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
