package handlers

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/templates"
	"net/http"
)

type GetIncidentsHandler struct{}

func NewGetIncidentsHandler() *GetIncidentsHandler {
	return &GetIncidentsHandler{}
}

func (h *GetIncidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	c := templates.Incidents(userInfo)

	err = templates.Layout("Incidents", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
