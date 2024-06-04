package handlers

import (
	"github.com/Ryan-Har/site-monitor/src/templates"
	"log/slog"
	"net/http"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Login()

	err := templates.Layout("Login", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering login template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
