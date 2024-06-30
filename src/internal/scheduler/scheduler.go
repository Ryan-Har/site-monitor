package scheduler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"github.com/Ryan-Har/site-monitor/src/internal/notifier"
	"github.com/Ryan-Har/site-monitor/src/internal/requests"
)

func StartSchedulers(dbh database.DBHandler) {
	resultsCh := make(chan []requests.Response, 10) // Buffered channel to hold results
	slog.Info("scheduler started")
	// Calculate the duration until the next minute mark
	now := time.Now()
	next := now.Truncate(time.Minute).Add(time.Minute)
	durationUntilNextMinute := time.Until(next)

	// Wait until the next minute mark
	time.Sleep(durationUntilNextMinute)

	// Start the ticker to trigger every minute on the minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	//run once after sleep, ticker handles after
	go processStart(dbh, resultsCh)

	// Process results asynchronously
	go processResults(dbh, resultsCh)

	for range ticker.C {
		go processStart(dbh, resultsCh)
	}
}

func getDivisorsOfMinute(currentMin int) []int {
	var result []int
	sqrtN := int(math.Sqrt(float64(currentMin)))
	for i := 1; i <= sqrtN; i++ {
		if currentMin%i == 0 {
			result = append(result, i)
			if i != currentMin/i {
				result = append(result, currentMin/i)
			}
		}
	}
	return result
}

func convertMinuteToSecs(min ...int) []int {
	var result []int
	for _, i := range min {
		result = append(result, i*60)
	}
	return result
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func monitorToRequestsStruct(monitors ...database.Monitor) []requests.Requests {
	var result []requests.Requests
	for _, m := range monitors {
		req := requests.Requests{
			URL:     m.URL,
			RType:   requests.RequestType(m.Type),
			Port:    m.Port,
			ID:      m.MonitorID,
			Timeout: time.Duration(m.TimeoutSecs) * time.Second,
		}
		result = append(result, req)
	}
	return result
}

func requestsResponseToResultsStruct(req ...requests.Response) []database.MonitorResult {
	var results []database.MonitorResult
	for _, r := range req {
		if r.Err != nil {
			slog.Error("error occured with result of request check", "original_request", r.Requests, "err", r.Err.Error())
		}
		result := database.MonitorResult{
			MonitorID:      r.ID,
			IsUp:           btoi(r.Up),
			ResponseTimeMs: int(r.ResponseTime.Milliseconds()),
			RunTimeEpoch:   int(r.RunTime.Unix()),
		}
		results = append(results, result)
	}
	return results
}

func processStart(dbh database.DBHandler, results chan<- []requests.Response) {
	currentMin := time.Now().Minute()
	currentMinDivisors := getDivisorsOfMinute(currentMin)
	secondsFilter := convertMinuteToSecs(currentMinDivisors...)
	monitorBatch, err := dbh.GetMonitors(database.ByIntervalSecs{Intervals: secondsFilter})
	if err != nil {
		slog.Error("error getting monitors by interval secs", "filter", secondsFilter, "err", err.Error())
	}

	reqBatch := monitorToRequestsStruct(monitorBatch...)
	result := requests.Send(reqBatch...)
	results <- result
}

func processResults(dbh database.DBHandler, results <-chan []requests.Response) {
	for msg := range results {
		dbResults := requestsResponseToResultsStruct(msg...)
		if err := dbh.AddMonitorResults(dbResults...); err != nil {
			slog.Error("error adding monitor results to db", "err", err.Error())
			continue
		}

		for _, response := range msg {
			dbResults := requestsResponseToResultsStruct(response)
			result := dbResults[0]

			chgType, err := getIncidentChangeType(dbh, result)
			if err != nil {
				slog.Info("error determining change type for result", "err", err, "result", result)
			}
			if chgType == changeTypeNoChange {
				continue
			}
			slog.Info("change type check", "changetype", chgType)
			notis, err := dbh.GetMonitorNotifications(database.ByMonitorIds{Ids: []int{result.MonitorID}})
			if err != nil {
				slog.Error("error getting monitor notifications for monitor id", "err", err.Error(), "monitorid", result.MonitorID)
			}
			slog.Info("notis check", "notis", notis)
			if chgType == changeTypeStart {
				incident := database.Incident{
					StartTime: result.RunTimeEpoch,
					MonitorID: result.MonitorID,
				}
				if err := dbh.AddIncidents(incident); err != nil {
					slog.Error("error adding incident to database", "err", err)
				}
			}
			if chgType == changeTypeEnd {
				resp, err := dbh.GetIncidents(database.ByMonitorIds{Ids: []int{result.MonitorID}}, database.IsOngoing{})
				if err != nil || len(resp) == 0 {
					slog.Error("error getting existing incident from database to close", "err", err, "checkid", result.CheckID)
					continue
				}

				//TODO: implement method to force db to respond with only a single item to make it more robust, although it should only ever respond with a single item if everything else is working as intended
				existingInc := resp[0]
				existingInc.EndTime = sql.NullInt64{Int64: int64(result.RunTimeEpoch), Valid: true}
				if err := dbh.CloseIncident(existingInc); err != nil {
					slog.Error("error closing incident", "err", err, "incident", existingInc)
				}
			}

			var notiIDs []int
			for _, noti := range notis {
				notiIDs = append(notiIDs, noti.NotificationID)
			}

			var notiSettings []database.NotificationSettings
			if len(notiIDs) > 0 {
				settings, err := dbh.GetNotifications(database.ByNotificationIds{Ids: notiIDs})
				if err != nil {
					slog.Error("error getting notifications for notification ids", "err", err.Error(), "notificationid", notiIDs)
				}
				notiSettings = settings
			}

			var upDown string
			if response.Up {
				upDown = "UP"
			} else {
				upDown = "DOWN"
			}
			for _, noti := range notiSettings {

				switch noti.NotificationType {
				case database.TypeDiscord:
					notifier.NewDiscordNotifier(
						notifier.WithUrl(noti.AdditionalInfo)).SendMessage(
						fmt.Sprintf("ALERT: Monitor %s: %s is now %s", response.RType, response.URL, upDown))
				case database.TypeSlack:
					notifier.NewSlackNotifier(
						notifier.WithUrl(noti.AdditionalInfo)).SendMessage(
						fmt.Sprintf("ALERT: Monitor %s: %s is now %s", response.RType, response.URL, upDown))
				}
			}

			// 	for _, result := range dbResults {
			// 		chgType, err := getIncidentChangeType(dbh, result)
			// 		if err != nil {
			// 			slog.Info("error determining change type for result", "err", err, "result", result)
			// 			continue
			// 		}
			// 		if chgType == changeTypeNoChange {
			// 			continue
			// 		}

			// 		notis, err := dbh.GetMonitorNotifications(database.ByMonitorIds{Ids: []int{result.MonitorID}})
			// 		if err != nil {
			// 			slog.Error("error getting monitor notifications for monitor id", "err", err.Error(), "monitorid", result.MonitorID)
			// 			continue
			// 		}

			// 		var notiIDs []int
			// 		for _, noti := range notis{
			// 			notiIDs = append(notiIDs, noti.NotificationID)
			// 		}

			// 		var notiSettings []database.NotificationSettings
			// 		if len(notiIDs) > 0 {
			// 			settings, err := dbh.GetNotifications(database.ByNotificationIds{Ids: notiIDs})
			// 			if err != nil {
			// 				slog.Error("error getting notifications for notification ids", "err", err.Error(), "notificationid", notiIDs)
			// 			}
			// 			notiSettings = settings
			// 		}

			// 		if chgType == changeTypeStart {
			// 			incident := database.Incident{
			// 				StartTime: result.RunTimeEpoch,
			// 				MonitorID: result.MonitorID,
			// 			}
			// 			if err := dbh.AddIncidents(incident); err != nil {
			// 				slog.Error("error adding incident to database", "err", err)
			// 			}

			// 			if len(notiSettings) < 1 {
			// 				continue
			// 			}

			// 			for _, noti := range notiSettings {
			// 				switch noti.NotificationType {
			// 				case database.TypeDiscord:
			// 					notifier.NewDiscordNotifier(
			// 						notifier.WithUrl(noti.AdditionalInfo)).SendMessage(
			// 							fmt.Sprintf("ALERT: Monitor ", result.))
			// 				case database.TypeSlack:
			// 				}
			// 			}

			// 		}
			// 		if chgType == changeTypeEnd {
			// 			resp, err := dbh.GetIncidents(database.ByMonitorIds{Ids: []int{result.MonitorID}}, database.IsOngoing{})
			// 			if err != nil || len(resp) == 0 {
			// 				slog.Error("error getting existing incident from database to close", "err", err, "checkid", result.CheckID)
			// 			}

			// 			//TODO: implement method to force db to respond with only a single item to make it more robust, although it should only ever respond with a single item if everything else is working as intended
			// 			existingInc := resp[0]
			// 			existingInc.EndTime = sql.NullInt64{Int64: int64(result.RunTimeEpoch), Valid: true}
			// 			if err := dbh.CloseIncident(existingInc); err != nil {
			// 				slog.Error("error closing incident", "err", err, "incident", existingInc)
			// 			}

			// 		}
		}
	}
}

type changeType string

const (
	changeTypeStart    = "start"
	changeTypeEnd      = "end"
	changeTypeNoChange = "noChange"
)

func getIncidentChangeType(dbh database.DBHandler, result database.MonitorResult) (changeType, error) {

	otherResults, err := dbh.GetMonitorResults(database.ByMonitorIds{Ids: []int{result.MonitorID}})
	if err != nil {
		return changeTypeNoChange, err
	}

	if len(otherResults) <= 1 {
		return changeTypeNoChange, nil
	}

	lastResult := otherResults[len(otherResults)-2]

	combined := result.IsUp ^ lastResult.IsUp
	noChange := combined&1 == 0

	if noChange {
		return changeTypeNoChange, nil
	}
	if lastResult.IsUp&1 == 1 {
		return changeTypeStart, nil
	}

	return changeTypeEnd, nil
}
