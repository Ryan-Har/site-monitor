package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ryan-Har/site-monitor/src/internal/database"

	"github.com/Ryan-Har/site-monitor/src/templates/partials"
)

type PostNotificationSettingsHandler struct {
	dbHandler database.DBHandler
}

func NewPostNotificationSettingsHandler(dbh database.DBHandler) *PostNotificationSettingsHandler {
	return &PostNotificationSettingsHandler{
		dbHandler: dbh,
	}
}

func (h PostNotificationSettingsHandler) ByID(w http.ResponseWriter, r *http.Request) {
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

	value := r.FormValue("additionalinfo")
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "value cannot be empty")
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

	updateStruct := database.NotificationSettings{
		Notificationid: id,
		AdditionalInfo: value,
	}

	err = h.dbHandler.UpdateNotificationAdditionalInfo(updateStruct)
	if err != nil {
		fmt.Fprintf(w, "error updating database, please try again")
		return
	}

	notificationResponse[0].AdditionalInfo = value

	err = partials.ExistingNotifications(notificationResponse[0]).Render(r.Context(), w)
	if err != nil {
		slog.Error("error rendering edit notification by id")
	}
}
