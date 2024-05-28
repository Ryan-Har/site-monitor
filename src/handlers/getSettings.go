package handlers

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/templates"
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
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
	nav := partials.SettingsNavBar("account")
	acc := templates.SettingsAccount("Ryan Harris")
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
	}

	c := templates.Settings(nav, acc, userInfo)
	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetNotificationSettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nav := partials.SettingsNavBar("notifications")
	not := templates.SettingsNotifications()
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	c := templates.Settings(nav, not, userInfo)

	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetSecuritySettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nav := partials.SettingsNavBar("security")
	sec := templates.SettingsSecurity()
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
	}

	c := templates.Settings(nav, sec, userInfo)
	err = templates.Layout("Settings", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
