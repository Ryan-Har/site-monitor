package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"github.com/Ryan-Har/site-monitor/src/models"
	"github.com/Ryan-Har/site-monitor/src/templates"
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
)

type GetMonitorOverviewHandler struct {
	dbHandler database.DBHandler
}

func NewGetMonitorOverviewHandler(dbh database.DBHandler) *GetMonitorOverviewHandler {
	return &GetMonitorOverviewHandler{
		dbHandler: dbh,
	}
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
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	userMonitors, err := h.dbHandler.GetMonitors(database.ByUUIDs{Ids: []string{userInfo.UUID}})
	if err != nil {
		fmt.Println(userMonitors)
		//handle error by sending no cards, not monitors are setup
	}

	usersMonitorIDs := make([]int, 0, len(userMonitors))
	userMonitorInfoMap := make(map[int]database.Monitor)
	for _, usermon := range userMonitors {
		usersMonitorIDs = append(usersMonitorIDs, usermon.MonitorID)
		userMonitorInfoMap[usermon.MonitorID] = usermon
	}
	monitorResults, err := h.dbHandler.GetMonitorResults(
		database.ByMonitorIds{Ids: usersMonitorIDs},
		database.BetweenRunTime{MinEpoch: unixTime7DaysAgo(), MaxEpoch: unixTimeNow()})
	if err != nil {
		slog.Error("error getting monitor results") //handle better
	}

	userMonitorResultsMap := make(map[int][]database.MonitorResult)

	for _, userRes := range monitorResults {
		userMonitorResultsMap[userRes.MonitorID] = append(userMonitorResultsMap[userRes.MonitorID], userRes)
	}

	cards := make([]models.MonitorCardGenerationModel, 0, len(usersMonitorIDs))
	for _, monitorId := range usersMonitorIDs {
		var currentMonitor models.MonitorCardGenerationModel

		name := fmt.Sprintf("%s: %s", userMonitorInfoMap[monitorId].Type, userMonitorInfoMap[monitorId].URL)
		interval := userMonitorInfoMap[monitorId].IntervalSecs

		currentMonitor.Name = name
		currentMonitor.RefreshIntervalSecs = interval

		// if not ok then no checks have happened yet
		monitorIdResults, ok := userMonitorResultsMap[monitorId]
		if ok {
			lastResult := monitorIdResults[len(monitorIdResults)-1]
			isLastUp := lastResult.IsUp == 1
			currentMonitor.Up = isLastUp
		}
		cards = append(cards, currentMonitor)
	}

	// for k, v := range userMonitorInfoMap {
	// 	slog.Info("infomap: ", "id", k, "info", v)
	// }

	// for k, v := range userMonitorResultsMap {
	// 	slog.Info("resultsmap: ", "id", k, "results", v)
	// }

	c := templates.MonitorOverview(userInfo, cards...)

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

func unixTimeNow() int {
	return int(time.Now().Unix())
}

func unixTime7DaysAgo() int {
	return int(time.Now().Add(-7 * 24 * time.Hour).Unix())
}
