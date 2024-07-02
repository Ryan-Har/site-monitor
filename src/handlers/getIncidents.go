package handlers

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"github.com/Ryan-Har/site-monitor/src/templates"
)

type GetIncidentsHandler struct {
	dbHandler database.DBHandler
}

func NewGetIncidentsHandler(dbh database.DBHandler) *GetIncidentsHandler {
	return &GetIncidentsHandler{
		dbHandler: dbh,
	}
}

func (h *GetIncidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Warn("error getting user info from context")
	}
	incWithMonitor, err := h.dbHandler.GetIncidentsWithMonitorInfoByUUID(userInfo.UUID)
	if err != nil {
		slog.Error("error getting incidents with monitor info by uuid", "err", err)
	}

	slices.Reverse(incWithMonitor) // reverse so newest if first
	c := templates.Incidents(userInfo, incWithMonitor...)
	err = templates.Layout("Incidents", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering incidents template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
