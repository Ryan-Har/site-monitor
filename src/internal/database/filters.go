package database

import (
	"fmt"
)

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
	placeholders := make([]interface{}, len(f.Intervals))
	for i, id := range f.Intervals {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Interval_in_seconds IN (%s) ", monitorTable, generateQuestionMarks(len(f.Intervals))), placeholders
}

type ByTimeoutSecs struct {
	Timeouts []int
}

func (f ByTimeoutSecs) MonitorToSQLite(monitorTable string) (string, []interface{}) {
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
	placeholders := make([]interface{}, len(f.Ports))
	for i, id := range f.Ports {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Port IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ports))), placeholders
}

type MonitorResultsFilter interface {
	ResultsToSQLite(monitorTable string) (string, []interface{})
}

type ByCheckIds struct {
	Ids []int
}

func (f ByCheckIds) ResultsToSQLite(monitorTable string) (string, []interface{}) {
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Check_id IN (%s) ", monitorTable, generateQuestionMarks(len(f.Ids))), placeholders
}

func (f ByMonitorIds) ResultsToSQLite(monitorTable string) (string, []interface{}) {
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
	MinEpoch int
	MaxEpoch int
}

func (f BetweenRunTime) ResultsToSQLite(monitorTable string) (string, []interface{}) {

	placeholder := make([]interface{}, 2)
	placeholder[0] = f.MinEpoch
	placeholder[1] = f.MaxEpoch

	return fmt.Sprintf(" %s.Run_time BETWEEN ? AND ? ", monitorTable), placeholder
}

type NotificationFilter interface {
	NotificationToSQLite(notificationTable string) (string, []interface{})
}

func (f ByUUIDs) NotificationToSQLite(notificationTable string) (string, []interface{}) {
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.UUID IN (%s) ", notificationTable, generateQuestionMarks(len(f.Ids))), placeholders
}

type ByNotificationIds struct {
	Ids []int
}

func (f ByNotificationIds) NotificationToSQLite(notificationTable string) (string, []interface{}) {
	placeholders := make([]interface{}, len(f.Ids))
	for i, id := range f.Ids {
		placeholders[i] = id
	}

	return fmt.Sprintf(" %s.Notification_id IN (%s) ", notificationTable, generateQuestionMarks(len(f.Ids))), placeholders
}
