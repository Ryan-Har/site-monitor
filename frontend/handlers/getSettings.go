package handlers

import (
	"github.com/Ryan-Har/site-monitor/frontend/templates"
	"github.com/Ryan-Har/site-monitor/frontend/templates/partials"
	"net/http"
)

type GetAccountSettingsHandler struct{}

func NewGetAccountSettingsHandler() *GetAccountSettingsHandler {
	return &GetAccountSettingsHandler{}
}

type GetNotificationSettingsHandler struct{}

func NewGetNotificationSettingsHandler() *GetNotificationSettingsHandler {
	return &GetNotificationSettingsHandler{}
}

type GetSecuritySettingsHandler struct{}

func NewGetSecuritySettingsHandler() *GetSecuritySettingsHandler {
	return &GetSecuritySettingsHandler{}
}

func (h *GetAccountSettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := partials.SettingsNavBar("account")
	c := templates.Settings(n)
	err := templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetNotificationSettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := partials.SettingsNavBar("notifications")
	c := templates.Settings(n)
	err := templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetSecuritySettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := partials.SettingsNavBar("security")
	c := templates.Settings(n)
	err := templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
