package handlers

import (
	"github.com/Ryan-Har/site-monitor/frontend/templates"
	"net/http"
)

type GetSignupHandler struct{}

func NewGetSignupHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetSignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Signup()

	err := templates.Layout("Sign up", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
