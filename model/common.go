package model

import (
	"encoding/json"
	"time"
)

// User : interface type of 'RealUser' and 'NullUser'
type User interface {
	ReportStatus(payload interface{})
}

// RelayState represents state of a relay (On/Off) and it's detail
type RelayState struct { // use when set relay mode
	Mode   string      `json:"mode" binding:"required"` // ON OFF AUTO SCHED ...
	Detail interface{} `json:"detail"`                  // detail depending on mode
}

// scheduleEntry represents user's schedule in TIMER mode
type scheduleEntry struct {
	StartHour  int   `json:"startHour"`
	StartMin   int   `json:"startMin"`
	EndHour    int   `json:"endHour"`
	EndMin     int   `json:"endMin"`
	DayOfWeeks []int `json:"dows"` // array of numbers in 0-6 represnting day of week this this schedule is active
}

// ScheduleDetail wraps schedule array
type ScheduleDetail struct {
	Schedules []scheduleEntry `json:"schedules" binding:"required"`
}

// DeviceSchema : basic device info
type DeviceSchema struct {
	ID          string                `json:"id"`
	Secret      string                `json:"secret"`
	Owner       string                `json:"owner"`
	RelayStates map[string]RelayState `json:"state"`
}

// DeviceMessage : mqtt stat from device
type DeviceMessage struct {
	Type string `json:"t" binding:"required"`
	// Payload struct {
	Soil     float32 `json:"Soil"`
	Humidity float32 `json:"Humidity"`
	Temp     float32 `json:"Temp"`
	// Also relay state ON/OFF ??
	// } `json:"data"`
}

// ToMap is convenient method for converting struct back to map
func (dmesg *DeviceMessage) ToMap() (out map[string]interface{}) {
	str, _ := json.Marshal(dmesg)
	json.Unmarshal(str, &out)
	return
}

// shortcut to create time for today with just HH:MM
func createTime(hour, min int) int64 {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location()).UnixNano()
}

// ToDeviceState convert time schedule to [][2]int if it's mode is TIMER
func (state *RelayState) ToDeviceState() RelayState {
	cpy := *state // copy it
	if state.Mode == "TIMER" {
		dow := time.Now().Weekday()

		var schedules ScheduleDetail
		str, _ := json.Marshal(state.Detail)
		json.Unmarshal(str, &schedules)

		schedArray := make([][2]int64, 0, 4)
		for _, sched := range schedules.Schedules {
			for _, d := range sched.DayOfWeeks {
				if time.Weekday(d) == dow {
					var cur [2]int64
					cur[0] = createTime(sched.StartHour, sched.StartMin)
					cur[1] = createTime(sched.EndHour, sched.EndMin)
					schedArray = append(schedArray, cur)
					break
				}
			}
		}
		cpy.Detail = schedArray
	}
	return cpy
}
