package handlers

import (
	"github.com/Ryan-Har/site-monitor/frontend/templates"
	//"github.com/Ryan-Har/site-monitor/frontend/templates/partials"
	"net/http"
)

type GetMaintenanceHandler struct{}

func NewGetMaintenanceHandler() *GetMaintenanceHandler {
	return &GetMaintenanceHandler{}
}

func (h *GetMaintenanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Maintenance()

	err := templates.Layout("Maintenance", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
