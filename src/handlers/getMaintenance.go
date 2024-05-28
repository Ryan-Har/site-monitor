package handlers

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/templates"
	"net/http"
)

type GetMaintenanceHandler struct{}

func NewGetMaintenanceHandler() *GetMaintenanceHandler {
	return &GetMaintenanceHandler{}
}

func (h *GetMaintenanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	c := templates.Maintenance(userInfo)

	err = templates.Layout("Maintenance", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
