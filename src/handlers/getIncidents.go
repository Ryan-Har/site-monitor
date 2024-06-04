package handlers

import (
	"github.com/Ryan-Har/site-monitor/src/templates"
	"log/slog"
	"net/http"
)

type GetIncidentsHandler struct{}

func NewGetIncidentsHandler() *GetIncidentsHandler {
	return &GetIncidentsHandler{}
}

func (h *GetIncidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Warn("error getting user info from context")
	}
	c := templates.Incidents(userInfo)

	err = templates.Layout("Incidents", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering incidents template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
