package models

type MonitorCardGenerationModel struct {
	MonitorID           int
	Up                  bool //is it currently up?
	Name                string
	RefreshIntervalSecs int //how often is the check run?
	LastChangeSecs      int //when the status last changed
	LastCheckSecs       int //how long ago was the last check, in seconds?
}
