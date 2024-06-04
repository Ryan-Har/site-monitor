package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Ryan-Har/site-monitor/src/templates"
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
)

type GetMonitorOverviewHandler struct{}

func NewGetMonitorOverviewHandler() *GetMonitorOverviewHandler {
	return &GetMonitorOverviewHandler{}
}

type GetMonitorFormHandler struct{}

func NewGetMonitorFormHandler() *GetMonitorFormHandler {
	return &GetMonitorFormHandler{}
}

type GetMonitorByID struct{}

func NewGetMonitorByID() *GetMonitorByID {
	return &GetMonitorByID{}
}

func (h *GetMonitorOverviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mCards := partials.MultipleMonitors()
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Warn("error getting user info from context")
	}
	c := templates.MonitorOverview(mCards, userInfo)

	err = templates.Layout("Monitors", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering incidents template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetMonitorFormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Warn("error getting user info from context")
	}
	c := templates.NewMonitorForm(userInfo)

	err = templates.Layout("Monitors", c).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *GetMonitorFormHandler) ServeFormContent(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	queryString := url.Query()

	typeSelection, ok := queryString["typeSelection"]
	if !ok {
		http.Error(w, "typeSelection not found in query string", http.StatusBadRequest)
		return
	}

	//should only ever include a single option, so we'll take the first one
	switch typeSelection[0] {
	case "HTTP":
		if err := partials.MonitorFormContentHTTP().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "ICMP":
		if err := partials.MonitorFormContentPing().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "TCP":
		if err := partials.MonitorFormContentPort().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid typeSelection", http.StatusBadRequest)
		return
	}
}

func (h *GetMonitorByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("monitorid")
	if id == "" {
		http.Error(w, "monitorid not found in query string", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Fprintf(w, "Getting monitor with id: %s", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
