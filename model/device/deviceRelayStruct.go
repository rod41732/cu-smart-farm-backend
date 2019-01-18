package device

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
)

// RelayState : state of relay (ON/OFF/auto + extra detail of mode)
type RelayState struct { // use when set relay mode
	Mode   string      `json:"mode"`   // ON OFF auto SCHED ...
	Detail interface{} `json:"detail"` // detail depending on mode
}

// scheduleEntry represents user's schedule in scheduled mode
type scheduleEntry struct {
	StartHour int `json:"startHour"`
	StartMin  int `json:"startMin"`
	EndHour   int `json:"endHour"`
	EndMin    int `json:"endMin"`
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
	Repeat    bool            `json:"repeat"`
	CreatedAt time.Time       `json:"createdAt"`
}

//FromMap is "constructor" for converting map[string]interface{} to Condition return error if can't convert
func (scheduleDetail *ScheduleDetail) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, scheduleDetail)
	if err != nil {
		return err
	}

	// Verify
	if scheduleDetail.CreatedAt.IsZero() {
		return errors.New("Empty time specified")
	}
	for _, entry := range scheduleDetail.Schedules {
		for _, h := range []int{entry.EndHour, entry.StartHour} {
			if !(0 <= h && h < 24) {
				return errors.New("Invalid Detail - Hour")
			}
		}
		for _, m := range []int{entry.EndMin, entry.StartMin} {
			if !(0 <= m && m < 60) {
				return errors.New("Invalid Detail - Min")
			}
		}
		if 60*entry.StartHour+entry.StartMin >= 60*entry.EndHour+entry.EndMin {
			return errors.New("Invalid Detail - Bad range")
		}
	}
	return nil
}

// shortcut to create time for today with just HH:MM
func createTime(hour, min int) int64 {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location()).Unix()
}

// ToDeviceState convert time schedule to [][2]int if it's mode is scheduled
func (state *RelayState) ToDeviceState() RelayState {
	cpy := *state // copy it
	if state.Mode == "scheduled" {
		var schedules ScheduleDetail
		str, _ := json.Marshal(state.Detail)
		json.Unmarshal(str, &schedules)
		detailStr := "off"
		now := time.Now().Unix()
		for _, sched := range schedules.Schedules {
			start := createTime(sched.StartHour, sched.StartMin)
			end := createTime(sched.EndHour, sched.EndMin)
			if start <= now && now <= end {
				detailStr = "on"
			}
		}
		cpy.Mode = "manual"
		cpy.Detail = detailStr
	}
	return cpy
}

// Verify verifys validity of RelayState object
func (state *RelayState) Verify() bool {
	if state.Mode == "manual" {
		str, ok := state.Detail.(string)
		return ok && common.StringInSlice(str, []string{"on", "off"})
	} else if state.Mode == "auto" {
		_map, ok := state.Detail.(map[string]interface{})
		if ok {
			var cond Condition
			err := cond.FromMap(_map)
			if err == nil {
				state.Detail = cond
				return true
			}
			return false
		}
		return false
	} else if state.Mode == "scheduled" {
		_map, ok := state.Detail.(map[string]interface{})
		if ok {
			var sched ScheduleDetail
			err := sched.FromMap(_map)
			if err == nil {
				state.Detail = sched
				return true
			}
			return false
		}
		return false
	}
	return false
}
