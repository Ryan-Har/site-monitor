package database

import (
	"errors"
	"fmt"
	"strings"
)

const (
	TypeEmail   = "email"
	TypeSlack   = "slack"
	TypeDiscord = "discord"
)

type NotificationType string

func (n NotificationType) String() string {
	return string(n)
}

type NotificationSettings struct {
	Notificationid   int
	UUID             string
	NotificationType NotificationType
	AdditionalInfo   string
}

func (h *SQLiteHandler) AddNotification(notificationSettings NotificationSettings) error {
	stmt, err := h.DB.Prepare("INSERT INTO Notifications (UUID, Type, Additional_info) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()

	_, err = stmt.Exec(notificationSettings.UUID, notificationSettings.NotificationType, notificationSettings.AdditionalInfo)
	if err != nil {
		h.writeMutex.Unlock()
		return err
	}

	h.writeMutex.Unlock()
	return nil
}

func (h *SQLiteHandler) GetNotifications(filters ...NotificationFilter) ([]NotificationSettings, error) {
	sqlTable := "Notifications"
	query := "SELECT Notification_id, UUID, Type, Additional_info FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.NotificationToSQLite(sqlTable)
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

	var notifications []NotificationSettings
	for rows.Next() {
		var n NotificationSettings
		err := rows.Scan(&n.Notificationid, &n.UUID, &n.NotificationType, &n.AdditionalInfo)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (h *SQLiteHandler) DeleteNotifications(filters ...NotificationFilter) error {
	if len(filters) == 0 {
		return fmt.Errorf("filters must accompany delete notifications")
	}

	sqlTable := "Notifications"
	query := "DELETE FROM " + sqlTable

	var whereClause []string
	var args []interface{}
	for _, filter := range filters {
		condSql, condArgs := filter.NotificationToSQLite(sqlTable)
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

func (h *SQLiteHandler) UpdateNotificationAdditionalInfo(notificationSettings NotificationSettings) error {
	if notificationSettings.Notificationid == 0 {
		return errors.New("notificationSettings struct does not include the notification id")
	}

	stmt, err := h.DB.Prepare("UPDATE Notifications SET Additional_info = ? WHERE Notification_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the prepared statement

	h.writeMutex.Lock()

	_, err = stmt.Exec(notificationSettings.AdditionalInfo, notificationSettings.Notificationid)
	if err != nil {
		h.writeMutex.Unlock()
		return err
	}

	h.writeMutex.Unlock()
	return nil
}
