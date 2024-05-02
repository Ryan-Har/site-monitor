package handlers

import (
	"github.com/Ryan-Har/site-monitor/frontend/templates"
	"github.com/Ryan-Har/site-monitor/frontend/templates/partials"
	"net/http"
)

type GetMonitorOverviewHandler struct{}

func NewGetMonitorOverviewHandler() *GetMonitorOverviewHandler {
	return &GetMonitorOverviewHandler{}
}

func (h *GetMonitorOverviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mCards := partials.MultipleMonitors()
	c := templates.MonitorOverview(mCards)

	err := templates.Layout("Monitors", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
