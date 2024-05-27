package models

type ExampleMonitor struct {
	Name          string
	URL           string
	Interval      int
	CurrentStatus bool
}

type MonitorStatus struct {
	MonitorName string
	Status      Status
}

type Status int

const (
	StatusUp Status = iota
	StatusDown
	StatusPaused
)

type MonitorStatusViewModel struct {
	Status map[Status]string
}
