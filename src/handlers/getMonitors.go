package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

type GetMonitorFormHandler struct {
	dbHandler database.DBHandler
}

func NewGetMonitorFormHandler(dbh database.DBHandler) *GetMonitorFormHandler {
	return &GetMonitorFormHandler{
		dbHandler: dbh,
	}
}

type GetMonitorByID struct {
	dbHandler database.DBHandler
}

func NewGetMonitorByID(dbh database.DBHandler) *GetMonitorByID {
	return &GetMonitorByID{
		dbHandler: dbh,
	}
}

type DeleteMonitorByID struct {
	dbHandler database.DBHandler
}

func NewDeleteMonitorByID(dbh database.DBHandler) *DeleteMonitorByID {
	return &DeleteMonitorByID{
		dbHandler: dbh,
	}
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
		var card models.MonitorCardGenerationModel

		// if not ok then no checks have happened yet
		monitorIdResults, ok := userMonitorResultsMap[monitorId]
		if ok {
			card = generateMonitorCardGenerationModel(userMonitorInfoMap[monitorId], monitorIdResults)
		} else {
			card.MonitorID = monitorId
			card.MType = userMonitorInfoMap[monitorId].Type
			card.MUrl = userMonitorInfoMap[monitorId].URL
			card.RefreshIntervalSecs = userMonitorInfoMap[monitorId].IntervalSecs
		}
		cards = append(cards, card)
	}

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

	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	resp, err := h.dbHandler.GetNotifications(database.ByUUIDs{Ids: []string{userInfo.UUID}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idNameMap := make(map[int]string)
	for _, v := range resp {
		idNameMap[v.Notificationid] = v.NotificationType.String()
	}

	//should only ever include a single option, so we'll take the first one
	switch typeSelection[0] {
	case "HTTP":
		if err := partials.MonitorFormContentHTTP(idNameMap).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "ICMP":
		if err := partials.MonitorFormContentPing(idNameMap).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "TCP":
		if err := partials.MonitorFormContentPort(idNameMap).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid typeSelection", http.StatusBadRequest)
		return
	}
}

func (h *GetMonitorByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for get monitor by id")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	idStr := r.PathValue("monitorid")
	if idStr == "" {
		http.Error(w, "monitorid not found in query string", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error converting monitor id to int")
		return
	}

	monitorOwnerCheckResponse, err := h.dbHandler.GetMonitors(database.ByMonitorIds{Ids: []int{id}})
	if err != nil || len(monitorOwnerCheckResponse) < 1 {
		fmt.Fprintf(w, "error getting monitor with id %d from database", id)
		return
	}

	if monitorOwnerCheckResponse[0].UUID != userInfo.UUID {
		http.Error(w, "monitor not owned by current user", http.StatusForbidden)
		return
	}

	checkResults, err := h.dbHandler.GetMonitorResults(database.ByMonitorIds{Ids: []int{id}})
	if err != nil {
		fmt.Fprintf(w, "error getting monitor results with id %d from database", id)
		return
	}

	monitorInfo := generateMonitorCardGenerationModel(monitorOwnerCheckResponse[0], checkResults)
	avg, min, max := getAvgMinMaxResponseMs(checkResults)
	responseTimeStats := partials.ResponseTimeStats(avg, min, max)

	c := templates.GetSingleMonitor(userInfo, monitorInfo, checkResults, responseTimeStats)

	err = templates.Layout("Monitor", c).Render(r.Context(), w)
	if err != nil {
		slog.Error("error while rendering single template", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *DeleteMonitorByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("monitorid")
	if idStr == "" {
		http.Error(w, "monitorid not found in query string", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error: unable to convert monitorid to int", http.StatusBadRequest)
		return
	}
	if err = h.dbHandler.DeleteMonitors(database.ByMonitorIds{Ids: []int{id}}); err != nil {
		http.Error(w, "error: unable to delete monitor", http.StatusInternalServerError)
		return
	}
	//return empty string, htmx is replacing the element with this
	fmt.Fprintf(w, "")
}

func unixTimeNow() int {
	return int(time.Now().Unix())
}

func unixTime7DaysAgo() int {
	return int(time.Now().Add(-7 * 24 * time.Hour).Unix())
}

func generateMonitorCardGenerationModel(monitor database.Monitor, monitorResult []database.MonitorResult) models.MonitorCardGenerationModel {
	var respCard models.MonitorCardGenerationModel

	respCard.MonitorID = monitor.MonitorID
	respCard.MType = monitor.Type
	respCard.MUrl = monitor.URL
	respCard.RefreshIntervalSecs = monitor.IntervalSecs
	if len(monitorResult) < 1 {
		return respCard
	}

	lastResult := monitorResult[len(monitorResult)-1]
	respCard.Up = lastResult.IsUp == 1
	for i := len(monitorResult) - 1; i >= 0; i-- {
		if monitorResult[i].IsUp != lastResult.IsUp {
			respCard.LastChangeSecs = lastResult.RunTimeEpoch - monitorResult[i].RunTimeEpoch
			break
		}
	}
	respCard.LastCheckSecs = unixTimeNow() - lastResult.RunTimeEpoch
	return respCard
}

type dateResponse struct {
	Date         string `json:"date"`
	ResponseTime int    `json:"responsetime"`
}

func (h *GetMonitorByID) ServeResponseTimes(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for get monitor by id")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	idStr := r.PathValue("monitorid")
	if idStr == "" {
		http.Error(w, "monitorid not found in query string", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error converting monitor id to int")
		return
	}

	monitorOwnerCheckResponse, err := h.dbHandler.GetMonitors(database.ByMonitorIds{Ids: []int{id}})
	if err != nil || len(monitorOwnerCheckResponse) < 1 {
		fmt.Fprintf(w, "error getting monitor with id %d from database", id)
		return
	}

	if monitorOwnerCheckResponse[0].UUID != userInfo.UUID {
		http.Error(w, "monitor not owned by current user", http.StatusForbidden)
		return
	}

	checkResults, err := h.dbHandler.GetMonitorResults(database.ByMonitorIds{Ids: []int{id}},
		database.BetweenRunTime{
			MinEpoch: unixTime7DaysAgo(),
			MaxEpoch: unixTimeNow(),
		})
	if err != nil {
		fmt.Fprintf(w, "error getting monitor results with id %d from database", id)
		return
	}

	dateResponseSlice := make([]dateResponse, 0, len(checkResults))

	for _, check := range checkResults {
		runTime64 := int64(check.RunTimeEpoch)
		d := dateResponse{
			Date:         time.Unix(runTime64, 0).Format("2006-01-02 15:04:05"),
			ResponseTime: check.ResponseTimeMs,
		}
		dateResponseSlice = append(dateResponseSlice, d)
	}

	jsonData, err := json.Marshal(dateResponseSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAvgMinMaxResponseMs(checks []database.MonitorResult) (avg int, min int, max int) {
	if len(checks) == 0 {
		return 0, 0, 0
	}
	var currentMin int = checks[0].ResponseTimeMs
	var currentMax int = checks[0].ResponseTimeMs
	var runningTotal int

	for _, mon := range checks {
		if mon.ResponseTimeMs < currentMin {
			currentMin = mon.ResponseTimeMs
		}
		if mon.ResponseTimeMs > currentMax {
			currentMax = mon.ResponseTimeMs
		}
		runningTotal += mon.ResponseTimeMs
	}
	mean := runningTotal / len(checks)
	return mean, currentMin, currentMax
}
