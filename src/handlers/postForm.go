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
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
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

	w.Header().Set("HX-Redirect", "/monitors")

	if err = r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	newMonitorForm := &newMonitorSetupForm{
		typeSelection:       r.Form.Get("typeSelection"),
		monitorLocation:     r.Form.Get("monitorLocation"),
		monitorIntervalMins: r.Form.Get("monitorIntervalNumber"),
		timeoutIntervalSecs: r.Form.Get("timeoutIntervalNumber"),
		portSelection:       r.Form.Get("monitorPort"),
	}

	monitor, err := newMonitorForm.validate()
	if err != nil {
		slog.Error("error validating new monitor form", "err", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "error validating new monitor form: %v", err)
		return
	}

	monitor.UUID = userInfo.UUID

	//default to 10 seconds for pings to prevent infinite waits later on
	if monitor.Type == "ICMP" && monitor.TimeoutSecs == 0 {
		monitor.TimeoutSecs = 10
	}

	if err = h.dbHandler.AddMonitors(*monitor); err != nil {
		slog.Error("error adding monitor to db", "monitor", monitor, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding monitor to db")
		return
	}

	notificationSelectionStr := r.Form.Get("notificationSelection")
	if notificationSelectionStr == "" {
		return
	}

	notificationSelection, err := strconv.Atoi(notificationSelectionStr)
	if err != nil {
		slog.Error("error converting notification ID to int")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error converting notification ID to integer type")
		return
	}

	newMonitor, err := h.dbHandler.GetMonitors(
		database.ByUUIDs{Ids: []string{monitor.UUID}},
		database.ByUrls{Urls: []string{monitor.URL}},
		database.ByTypes{Types: []string{monitor.Type}},
		database.ByIntervalSecs{Intervals: []int{monitor.IntervalSecs}},
		database.ByTimeoutSecs{Timeouts: []int{monitor.TimeoutSecs}},
		database.ByPorts{Ports: []int{monitor.Port}},
	)
	if err != nil || len(newMonitor) < 1 {
		slog.Error("error getting new monitor from db", "monitor", monitor, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding monitor to db")
		return
	}

	monitorToNotification := make(map[int]int)
	monitorToNotification[newMonitor[len(newMonitor)-1].MonitorID] = notificationSelection

	if err = h.dbHandler.AddMonitorNotification(monitorToNotification); err != nil {
		slog.Error("error adding new monitor notification", "monitor notification map", monitorToNotification, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding new monitor notification to db")
		return
	}

	w.Header().Set("HX-Redirect", "/monitors")
}

type newMonitorSetupForm struct {
	typeSelection       string
	monitorLocation     string
	monitorIntervalMins string // should be int
	timeoutIntervalSecs string // should be int
	portSelection       string // should be int
}

func (f *newMonitorSetupForm) validate() (*database.Monitor, error) {
	var validationErrors = map[string]string{
		"typeSelection":       "type Selection not one of: \"ICMP, TCP, HTTP\"",
		"monitorLocation":     "unable to parse monitorLocation",
		"monitorIntervalMins": "error converting monitorIntervalNumber to int type",
		"timeoutIntervalSecs": "error converting timeoutIntervalNumber to int type or value must be higher than 0",
		"portSelection":       "error converting portSelection to int type",
	}

	var monitorResponse database.Monitor

	if !map[string]bool{"ICMP": true, "TCP": true, "HTTP": true}[f.typeSelection] {
		return &database.Monitor{}, errors.New(validationErrors["typeSelection"])
	} else {
		monitorResponse.Type = f.typeSelection
	}

	if mins, err := stringToInt(f.monitorIntervalMins); err != nil {
		return &database.Monitor{}, errors.New(validationErrors["monitorIntervalMins"])
	} else {
		monitorResponse.IntervalSecs = mins * 60
	}

	switch f.typeSelection {
	case "HTTP":
		if !isURL(f.monitorLocation) {
			return &database.Monitor{}, errors.New(validationErrors["monitorLocation"])
		} else {
			monitorResponse.URL = f.monitorLocation
		}
		if secs, err := stringToInt(f.timeoutIntervalSecs); err != nil {
			return &database.Monitor{}, errors.New(validationErrors["timeoutIntervalSecs"])
		} else {
			monitorResponse.TimeoutSecs = secs
		}
	case "ICMP":
		if !(isDomain(f.monitorLocation) || isIPAddress(f.monitorLocation)) {
			return &database.Monitor{}, errors.New(validationErrors["monitorLocation"])
		} else {
			monitorResponse.URL = f.monitorLocation
		}
	case "TCP":
		if !(isDomain(f.monitorLocation) || isIPAddress(f.monitorLocation)) {
			return &database.Monitor{}, errors.New(validationErrors["monitorLocation"])
		} else {
			monitorResponse.URL = f.monitorLocation
		}
		if secs, err := stringToInt(f.timeoutIntervalSecs); err != nil {
			return &database.Monitor{}, errors.New(validationErrors["timeoutIntervalSecs"])
		} else {
			monitorResponse.TimeoutSecs = secs
		}
		if port, err := stringToInt(f.portSelection); err != nil {
			return &database.Monitor{}, errors.New(validationErrors["portSelection"])
		} else {
			monitorResponse.Port = port
		}
	}
	return &monitorResponse, nil
}

type ValidationFormHandler struct{}

func NewValidationFormHandler() *ValidationFormHandler {
	return &ValidationFormHandler{}
}

func (h ValidationFormHandler) ValidateMonitorLocationHttp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	errMsg := "input not valid, must be url"

	monitorLocation := r.Form.Get("monitorLocation")
	if isURL(monitorLocation) {
		err := partials.MonitorLocationHttpValidationResponseValid(monitorLocation).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location http valid response", "err", err.Error())
		}
	} else {
		err := partials.MonitorLocationHttpValidationResponseInvalid(monitorLocation, errMsg).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location http invalid response", "err", err.Error())
		}
	}
}

func (h ValidationFormHandler) ValidateMonitorLocationIpOrHost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	errMsg := "input not valid, must be domain or ip address"

	monitorLocation := r.Form.Get("monitorLocation")
	if isDomain(monitorLocation) || isIPAddress(monitorLocation) {
		err := partials.MonitorLocationIpOrHostValidationResponseValid(monitorLocation).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location ip or host valid response", "err", err.Error())
		}
	} else {
		err := partials.MonitorLocationIpOrHostValidationResponseInvalid(monitorLocation, errMsg).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location ip or host invalid response", "err", err.Error())
		}
	}
}

func (h ValidationFormHandler) ValidateMonitorPortNumber(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("error parsing form for newmonitor form post")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	errMsg := "port must be between 1 and 65535"

	portString := r.Form.Get("monitorPort")

	//edge case if user manually clears field, otherwise nothing is returned
	portInt, err := strconv.Atoi(portString)
	if err != nil {
		slog.Info("error validating monitor port number, not int, returning blank", "input", portString)
		err := partials.MonitorPortNumberValidationResponseInvalid(portString, errMsg).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location port number invalid response", "err", err.Error())
		}
		return
	}

	if isValidPortNumber(portInt) {
		err := partials.MonitorPortNumberValidationResponseValid(portString).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location port number valid response", "err", err.Error())
		}
	} else {
		err := partials.MonitorPortNumberValidationResponseInvalid(portString, errMsg).Render(r.Context(), w)
		if err != nil {
			slog.Error("error while rendering monitor location port number invalid response", "err", err.Error())
		}
	}
}

func stringToInt(val string) (int, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int, string provided: %s", val)
	}
	return i, nil
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

func isValidPortNumber(n int) bool {
	return 0 < n && n <= 65535
}

func (h *PostFormHandler) NewNotificationForm(w http.ResponseWriter, r *http.Request) {
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

	newNotification := &database.NotificationSettings{
		UUID:             userInfo.UUID,
		NotificationType: database.NotificationType(r.Form.Get("typeSelection")),
		AdditionalInfo:   r.Form.Get("additionalInfo"),
	}

	if err = h.dbHandler.AddNotification(*newNotification); err != nil {
		slog.Error("error adding notification method to db", "notification", newNotification, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding notification method to db")
		return
	}

	w.Header().Set("HX-Redirect", "/settings/notifications")
}
