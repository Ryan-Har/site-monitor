package database

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var testDBFile = "testdb.db"

func setupTesting() (*SQLiteHandler, func()) {
	db, err := openSQLiteDB(testDBFile)
	if err != nil {
		panic("error opening db: " + err.Error())
	}
	sqlite := SQLiteHandler{
		DB:         db,
		Version:    1,
		writeMutex: &sync.Mutex{},
	}
	cleanup := func() {
		os.Remove(testDBFile)
	}
	return &sqlite, cleanup
}

func getSampleMonitors(length int) []Monitor {
	var monitors []Monitor

	for l := range length {
		m := Monitor{
			MonitorID:    l + 1,
			UUID:         fmt.Sprintf("100%d", l+1),
			URL:          "https://example.com",
			Type:         "HTTP",
			IntervalSecs: 5,
			TimeoutSecs:  5,
			Port:         80,
		}
		monitors = append(monitors, m)
	}

	return monitors
}

func TestSQLiteFunctions(t *testing.T) {
	h, cleanup := setupTesting()
	defer cleanup()

	t.Run("Initialisation", func(t *testing.T) {
		if isSQLiteDBPopulated(h.DB) {
			t.Errorf("db is populated and shouldn't be. Is the testdb already populated?")
		}
		if err := populateSQLiteDB(h.DB); err != nil {
			t.Errorf("issue populating db with schema")
		}
		if !isSQLiteDBPopulated(h.DB) {
			t.Errorf("db should be populated and is not")
		}
	})

	t.Run("Add Single Monitor", func(t *testing.T) {
		m := getSampleMonitors(1)[0]
		err := h.AddMonitors(m)
		if err != nil {
			t.Errorf("error adding monitor, %v", err.Error())
		}
	})

	t.Run("Get Single Monitor by each filter", func(t *testing.T) {
		m := getSampleMonitors(1)[0]

		_, err := h.GetMonitors(ByMonitorIds{[]int{m.MonitorID}})
		if err != nil {
			t.Errorf("Get Single monitor by id filter failed, no results: %v", err.Error())
		}

		_, err = h.GetMonitors(ByUUIDs{[]string{m.UUID}})
		if err != nil {
			t.Errorf("Get Single monitor by UUID filter failed, no results: %v", err.Error())
		}

		_, err = h.GetMonitors(ByUrls{[]string{m.URL}})
		if err != nil {
			t.Errorf("Get Single monitor by URL filter failed, no results: %v", err.Error())
		}

		_, err = h.GetMonitors(ByTypes{[]string{m.Type}})
		if err != nil {
			t.Errorf("Get Single monitor by Type filter failed, no results: %v", err.Error())
		}

		_, err = h.GetMonitors(ByIntervalSecs{[]int{m.IntervalSecs}})
		if err != nil {
			t.Errorf("Get Single monitor by IntervalSecs filter failed, no results: %v", err.Error())
		}

		_, err = h.GetMonitors(ByTimeoutSecs{[]int{m.TimeoutSecs}})
		if err != nil {
			t.Errorf("Get Single monitor by TimeoutSecs filter failed, no results: %v", err.Error())
		}

		mon, err := h.GetMonitors(ByPorts{[]int{m.Port}})
		if err != nil {
			t.Errorf("Get Single monitor by Port filter failed, no results: %v", err.Error())
		}

		if len(mon) > 0 && mon[0] != m {
			t.Errorf("Monitor results not accurate")
		}
	})

	t.Run("Delete Single Monitor", func(t *testing.T) {
		m := getSampleMonitors(1)[0]
		monitorIds := []int{m.MonitorID}
		if err := h.DeleteMonitors(); err == nil {
			t.Errorf("Delete Monitor should have responded with error for no arguments, error was nil")
		}

		if err := h.DeleteMonitors(ByMonitorIds{monitorIds}); err != nil {
			t.Errorf("Delete Monitor failed: %v", err.Error())
		}
	})

	// t.Run("Add Multiple Monitors", func(t *testing.T) {

	// })

	// t.Run("Get Multiple Monitors", func(t *testing.T) {

	// })

	// t.Run("Delete Multiple Monitors", func(t *testing.T) {

	// })
}
