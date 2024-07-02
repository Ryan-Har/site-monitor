package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Incident struct {
	IncidentID int
	StartTime  int
	EndTime    sql.NullInt64
	MonitorID  int
}

func (h *SQLiteHandler) AddIncidents(incident ...Incident) error {
	stmt, err := h.DB.Prepare("INSERT INTO Incidents (Start_time, Monitor_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()
	for _, m := range incident {
		_, err := stmt.Exec(m.StartTime, m.MonitorID)
		if err != nil {
			h.writeMutex.Unlock()
			return err
		}
	}
	h.writeMutex.Unlock()
	return nil
}

func (h *SQLiteHandler) CloseIncident(incident Incident) error {
	if incident.IncidentID == 0 {
		return errors.New("incident struct does not include the incident id")
	}

	stmt, err := h.DB.Prepare("UPDATE Incidents SET End_time = ? WHERE Incident_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()
	defer h.writeMutex.Unlock()

	_, err = stmt.Exec(incident.EndTime.Int64, incident.IncidentID)
	if err != nil {
		return err
	}

	return nil
}

func (h *SQLiteHandler) GetIncidents(filters ...IncidentFilter) ([]Incident, error) {
	sqlTable := "Incidents"
	query := "SELECT Incident_id, Start_time, End_time, Monitor_id FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.IncidentToSQLite(sqlTable)
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

	var incidents []Incident
	for rows.Next() {
		var i Incident
		err := rows.Scan(&i.IncidentID, &i.StartTime, &i.EndTime, &i.MonitorID)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}
	return incidents, nil
}

func (h *SQLiteHandler) DeleteIncidents(filters ...IncidentFilter) error {
	if len(filters) == 0 {
		return fmt.Errorf("filters must accompany delete incidents")
	}

	sqlTable := "Incidents"
	query := "DELETE FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.IncidentToSQLite(sqlTable)
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

type IncidentWithMonitor struct {
	Incident
	UUID string
	URL  string
	Type string
	Port int
}

func (h *SQLiteHandler) GetIncidentsWithMonitorInfoByUUID(uuid string) ([]IncidentWithMonitor, error) {
	query := `SELECT i.Incident_id, i.Start_time, i.End_time, i.Monitor_id, m.UUID, m.Url, m.Type, m.Port 
	FROM Incidents i 
	LEFT JOIN Monitors m 
	ON i.Monitor_id = m.Monitor_id 
	WHERE m.UUID = ?`

	args := make([]interface{}, 1)
	args[0] = uuid

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []IncidentWithMonitor
	for rows.Next() {
		var i IncidentWithMonitor
		err := rows.Scan(&i.IncidentID, &i.StartTime, &i.EndTime, &i.MonitorID, &i.UUID, &i.URL, &i.Type, &i.Port)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}
	return incidents, nil
}
