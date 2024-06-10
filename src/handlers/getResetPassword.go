package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Ryan-Har/site-monitor/src/templates"
)

type GetResetPasswordHandler struct{}

func NewGetResetPasswordHandler() *GetResetPasswordHandler {
	return &GetResetPasswordHandler{}
}

func (h *GetResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.ResetPassword()

	err := templates.Layout("Reset Password", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering reset password template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
