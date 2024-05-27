package handlers

import (
	"github.com/Ryan-Har/site-monitor/src/templates"
	//"github.com/Ryan-Har/site-monitor/src/templates/partials"
	"net/http"
)

type GetIncidentsHandler struct{}

func NewGetIncidentsHandler() *GetIncidentsHandler {
	return &GetIncidentsHandler{}
}

func (h *GetIncidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Incidents()

	err := templates.Layout("Incidents", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
