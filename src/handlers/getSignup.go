package handlers

import (
	"github.com/Ryan-Har/site-monitor/src/templates"
	"log/slog"
	"net/http"
)

type GetSignupHandler struct{}

func NewGetSignupHandler() *GetSignupHandler {
	return &GetSignupHandler{}
}

func (h *GetSignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Signup()

	err := templates.Layout("Sign up", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering signup template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
