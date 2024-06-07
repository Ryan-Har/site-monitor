package handlers

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/templates"
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
	"log/slog"
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
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	nav := partials.SettingsNavBar("account")
	acc := templates.SettingsAccount(userInfo.Name)

	c := templates.Settings(nav, acc, userInfo)
	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering account settings template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetNotificationSettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	nav := partials.SettingsNavBar("notifications")
	not := templates.SettingsNotifications()

	c := templates.Settings(nav, not, userInfo)

	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering notifications settings template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetSecuritySettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	nav := partials.SettingsNavBar("security")
	sec := templates.SettingsSecurity()

	c := templates.Settings(nav, sec, userInfo)
	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering security settings template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
