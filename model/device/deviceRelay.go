package device

import (
	"encoding/json"
	"time"
)

// RelayState : state of relay (ON/OFF/AUTO + extra detail of mode)
type RelayState struct { // use when set relay mode
	Mode   string      `json:"mode" binding:"required"` // ON OFF AUTO SCHED ...
	Detail interface{} `json:"detail"`                  // detail depending on mode
}

// scheduleEntry represents user's schedule in TIMER mode
type scheduleEntry struct {
	StartHour  int   `json:"startHour" binding:"required"`
	StartMin   int   `json:"startMin" binding:"required"`
	EndHour    int   `json:"endHour" binding:"required"`
	EndMin     int   `json:"endMin" binding:"required"`
	DayOfWeeks []int `json:"dows" binding:"required"` // array of numbers in 0-6 represnting day of week this this schedule is active
}

// ScheduleDetail wraps schedule array
type ScheduleDetail struct {
	Schedules []scheduleEntry `json:"schedules" binding:"required"`
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

// Verify verifys validity of RelayState object
func (state *RelayState) Verify() bool {
	return true
}
