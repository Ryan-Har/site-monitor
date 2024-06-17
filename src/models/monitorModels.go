package models

import (
	"fmt"
)

type MonitorCardGenerationModel struct {
	MonitorID           int
	Up                  bool   //is it currently up?
	MType               string //monitor type
	MUrl                string //monitor url
	RefreshIntervalSecs int    //how often is the check run?
	LastChangeSecs      int    //when the status last changed
	LastCheckSecs       int    //how long ago was the last check, in seconds?
}

func (m *MonitorCardGenerationModel) Name() string {
	return fmt.Sprintf("%s: %s", m.MType, m.MUrl)
}
