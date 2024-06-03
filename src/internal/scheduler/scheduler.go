package scheduler

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"github.com/Ryan-Har/site-monitor/src/internal/requests"
	"math"
	"time"
)

func StartSchedulers(dbh database.DBHandler) {
	resultsCh := make(chan []requests.Response, 10) // Buffered channel to hold results
	fmt.Println("scheduler started")
	// Calculate the duration until the next minute mark
	now := time.Now()
	next := now.Truncate(time.Minute).Add(time.Minute)
	durationUntilNextMinute := time.Until(next)

	// Wait until the next minute mark
	time.Sleep(durationUntilNextMinute)

	// Start the ticker to trigger every minute on the minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	fmt.Println("first processStart run")
	//run once after sleep, ticker handles after
	go processStart(dbh, resultsCh)

	// Process results asynchronously
	go processResults(dbh, resultsCh)

	for range ticker.C {
		fmt.Println("ticker run, processing")
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
			fmt.Printf("error occured with result of request check. original request: %v, error: %s", r.Requests, r.Err.Error())
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
		fmt.Printf("error getting monitors by interval secs. filter: %v, err: %s \n", secondsFilter, err.Error())
	}

	reqBatch := monitorToRequestsStruct(monitorBatch...)
	result := requests.Send(reqBatch...)
	results <- result
}

func processResults(dbh database.DBHandler, results <-chan []requests.Response) {
	for msg := range results {
		dbResults := requestsResponseToResultsStruct(msg...)
		if err := dbh.AddMonitorResults(dbResults...); err != nil {
			fmt.Println("error adding monitor results to db: ", err.Error())
		}
	}
}
