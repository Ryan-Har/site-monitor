package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ryan-Har/site-monitor/src/internal/database"
)

type DeleteNotificationSettingsHandler struct {
	dbHandler database.DBHandler
}

func NewDeleteNotificationSettingsHandler(dbh database.DBHandler) *DeleteNotificationSettingsHandler {
	return &DeleteNotificationSettingsHandler{
		dbHandler: dbh,
	}
}

func (h DeleteNotificationSettingsHandler) ByID(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for get notification by id")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	idStr := r.PathValue("notificationid")
	if idStr == "" {
		http.Error(w, "notificationid not found in query string", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error converting notification id to int")
		return
	}

	notificationResponse, err := h.dbHandler.GetNotifications(database.ByNotificationIds{Ids: []int{id}})
	if err != nil || len(notificationResponse) < 1 {
		fmt.Fprintf(w, "error getting notification with id %d from database", id)
		return
	}

	if notificationResponse[0].UUID != userInfo.UUID {
		http.Error(w, "monitor not owned by current user", http.StatusForbidden)
		return
	}

	err = h.dbHandler.DeleteNotifications(database.ByNotificationIds{Ids: []int{id}})
	if err != nil {
		fmt.Fprintf(w, "error deleting database, please try again")
		return
	}

	fmt.Fprintf(w, "")
}
