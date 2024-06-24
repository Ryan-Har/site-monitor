package handlers

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/internal/notifier"
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

func (h *GetNotificationSettingsHandler) ServeFormContent(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryString := url.Query()

	typeSelection, ok := queryString["typeSelection"]
	if !ok {
		http.Error(w, "typeSelection not found in query string", http.StatusBadRequest)
		return
	}
	//should only ever include a single option, so we'll take the first one
	switch typeSelection[0] {
	case "discord", "slack":
		if err := partials.NotificationFormContentWebhook().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "email":
		if err := partials.NotificationFormContentEmail().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid typeSelection", http.StatusBadRequest)
		return
	}

}

func (h *GetNotificationSettingsHandler) SendTestNotification(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	typeSelection := r.Form.Get("typeSelection")
	endpoint := r.Form.Get("additionalInfo")

	if typeSelection == "" || endpoint == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "typeSelection or endpoint not populated")
		return
	}

	switch typeSelection {
	case "discord":
		err := notifier.NewDiscordNotifier(notifier.WithUrl(endpoint)).SendTest()
		if err != nil {
			fmt.Fprintf(w, "error sending discord notification %v", err)
			return
		} else {
			fmt.Fprintf(w, "discord notification sent successfully, please submit to save")
			return
		}
	case "slack":
		err := notifier.NewSlackNotifier(notifier.WithUrl(endpoint)).SendTest()
		if err != nil {
			fmt.Fprintf(w, "error sending slack notification %v", err)
			return
		} else {
			fmt.Fprintf(w, "slack notification sent successfully, please submit to save")
			return
		}
	case "email":
		fmt.Fprintf(w, "email notification sent successfully")
		return
	}

}
