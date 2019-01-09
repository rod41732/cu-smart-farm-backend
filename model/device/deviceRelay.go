package device

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
)

// RelayState : state of relay (ON/OFF/AUTO + extra detail of mode)
type RelayState struct { // use when set relay mode
	Mode   string      `json:"mode"`   // ON OFF AUTO SCHED ...
	Detail interface{} `json:"detail"` // detail depending on mode
}

// scheduleEntry represents user's schedule in TIMER mode
type scheduleEntry struct {
	StartHour  int   `json:"startHour"`
	StartMin   int   `json:"startMin"`
	EndHour    int   `json:"endHour"`
	EndMin     int   `json:"endMin"`
	DayOfWeeks []int `json:"dows"` // array of numbers in 0-6 represnting day of week this this schedule is active
}

// Condition is structure for condition in `auto` mode
type Condition struct {
	Sensor  string  `json:"sensor"`
	Trigger float32 `json:"trigger"`
	Symbol  string  `json:"symbol"`
}

// Validate validate value of Condition
func (condition *Condition) Validate() bool {
	if !(common.StringInSlice(condition.Sensor, []string{"soil", "temp", "humidity"})) {
		return false
	}
	if !(common.StringInSlice(condition.Symbol, []string{"<", ">"})) {
		return false
	}
	return true
}

//FromMap is "constructor" for converting map[string]interface{} to Condition return error if can't convert
func (condition *Condition) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, condition)
	if err != nil {
		return err
	}
	if condition.Validate() {
		return nil
	} else {
		return errors.New("Validation Error")
	}
}

// ScheduleDetail wraps schedule array
type ScheduleDetail struct {
	Schedules []scheduleEntry `json:"schedules"`
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
		if len(schedArray) == 0 && len(schedules.Schedules) != 0 {
			common.Println("[DEBUG] device state converter: there's data in schedule but day doesn't match")
		}
		cpy.Detail = schedArray
	}
	return cpy
}

// Verify verifys validity of RelayState object
func (state *RelayState) Verify() bool {
	return true
}
