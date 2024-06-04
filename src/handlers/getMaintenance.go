package handlers

import (
	"github.com/Ryan-Har/site-monitor/src/templates"
	"log/slog"
	"net/http"
)

type GetMaintenanceHandler struct{}

func NewGetMaintenanceHandler() *GetMaintenanceHandler {
	return &GetMaintenanceHandler{}
}

func (h *GetMaintenanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Warn("error getting user info from context")
	}
	c := templates.Maintenance(userInfo)

	err = templates.Layout("Maintenance", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering maintenance template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
