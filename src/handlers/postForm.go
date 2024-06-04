package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/Ryan-Har/site-monitor/src/internal/database"
)

type PostFormHandler struct {
	dbHandler database.DBHandler
}

func NewPostFormHandler(dbh database.DBHandler) *PostFormHandler {
	return &PostFormHandler{
		dbHandler: dbh,
	}
}

func (h *PostFormHandler) NewMonitorForm(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetUserInfoFromContext(r.Context())
	if err != nil {
		slog.Error("error getting user info from context for newmonitor form post")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error getting user info from context, reauthentication needed")
		return
	}

	if err = r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	monitor, err := validateNewMonitorForm(r.Form)
	if err != nil {
		slog.Error("error validating new monitor form", "err", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "error validating new monitor form: %v", err)
		return
	}

	monitor.UUID = userInfo.UUID

	if err = h.dbHandler.AddMonitors(monitor); err != nil {
		slog.Error("error adding monitor to db", "monitor", monitor, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding monitor to db")
		return
	}

}

func validateNewMonitorForm(v url.Values) (database.Monitor, error) {
	resp := database.Monitor{}

	monitorIntervalMinsString := v.Get("monitorIntervalNumber")
	timeoutIntervalSecsString := v.Get("timeoutIntervalNumber")
	typeSelection := v.Get("typeSelection")
	portSelectionStr := v.Get("monitorPort")
	url := v.Get("monitorLocation")

	acceptableTypes := map[string]bool{
		"ICMP": true,
		"TCP":  true,
		"HTTP": true,
	}
	if !acceptableTypes[typeSelection] {
		return database.Monitor{}, errors.New(`type Selection not one of: "ICMP, TCP, HTTP"`)
	}
	resp.Type = typeSelection

	switch typeSelection {
	case "HTTP":
		if !isURL(url) {
			return database.Monitor{}, errors.New("unable to parse monitorLocation as URL")
		}
	case "ICMP", "TCP":
		if !isDomain(url) && !isIPAddress(url) {
			return database.Monitor{}, errors.New("unable to parse monitorLocation as domain name or ip address")
		}
	default:
		return database.Monitor{}, errors.New("unable to validate monitorLocation")
	}
	resp.URL = url

	monitorIntervalMinsInt, err := strconv.Atoi(monitorIntervalMinsString)
	if err != nil {
		return database.Monitor{}, errors.New("error converting monitorIntervalNumber to int type")
	}
	monitorIntervalSeconds := monitorIntervalMinsInt * 60
	resp.IntervalSecs = monitorIntervalSeconds

	timeoutIntervalSecsInt, err := strconv.Atoi(timeoutIntervalSecsString)
	if err != nil {
		return database.Monitor{}, errors.New("error converting timeoutIntervalNumber to int type")
	}
	if timeoutIntervalSecsInt <= 0 {
		return database.Monitor{}, errors.New("error converting timeoutIntervalNumber must be higher than 0")
	}
	resp.TimeoutSecs = timeoutIntervalSecsInt

	if typeSelection == "TCP" {
		portSelectionInt, err := strconv.Atoi(portSelectionStr)
		if err != nil {
			return database.Monitor{}, errors.New("error converting portSelection to int type")
		}
		resp.Port = portSelectionInt
	}

	return resp, nil
}

func isDomain(domain string) bool {
	re := regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)
	return re.MatchString(domain)
}

// isURL checks if a string is a valid URL
func isURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

// isIPAddress checks if a string is a valid IP address (IPv4 or IPv6)
func isIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}
