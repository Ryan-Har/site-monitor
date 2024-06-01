package database

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var testDBFile = "testdb.db"

func setupTesting() (*SQLiteHandler, func()) {
	db, _ := openSQLiteDB(testDBFile)
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
		err := h.AddMonitor(m)
		if err != nil {
			t.Errorf("error adding monitor, %v", err.Error())
		}
	})

	t.Run("Get Single Monitor", func(t *testing.T) {
		m := getSampleMonitors(1)[0]

		if _, err := h.GetMonitorByID(); err == nil {
			t.Errorf("GetMonitorByID should have responded with error for no arguments, error was nil")
		}

		mon, err := h.GetMonitorByID(1)
		if err != nil {
			t.Errorf("Get Single monitor failed, no results: %v", err.Error())
		}

		if len(mon) > 0 && mon[0] != m {
			t.Errorf("Monitor results not accurate")
		}
	})

	t.Run("Delete Single Monitor", func(t *testing.T) {
		if err := h.DeleteMonitorByID(); err == nil {
			t.Errorf("Delete Monitor should have responded with error for no arguments, error was nil")
		}

		if err := h.DeleteMonitorByID(1); err != nil {
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
