package message

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model/device"
)

// APICall represent payload for API Calls via WS
type APICall struct {
	EndPoint string                 `json:"endPoint"` // addDevice, removeDevice, setDevice, pollDevice, listDevice...
	Token    string                 `json:"token"`
	Payload  map[string]interface{} `json:"payload"` // json data depend on command
}

// Message is regular message format for any device API
type Message struct {
	DeviceID string                 `json:"deviceID"`
	Param    map[string]interface{} `json:"param"` // json data depend on command
}

// AddDeviceMessage is payload format for addDevice API
type AddDeviceMessage struct {
	DeviceSecret string `json:"deviceSecret"`
	DeviceName   string `json:"deviceName"`
}

// DeviceCommandMessage is payload format for setDevice API
type DeviceCommandMessage struct {
	RelayID string            `json:"relayID"`
	State   device.RelayState `json:"state"`
}

// RenameDeviceMessage is payload format for rename device API
type RenameDeviceMessage struct {
	Name string `json:"name"`
}

// TimeQuery is payload for querying sersnor logs
type TimeQuery struct {
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Limit int       `json:"limit"`
}

// FromMap is "constructor" for converting map[string]interface{} to AddDeviceMessage  return error if can't convert
func (message *AddDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

// FromMap is "constructor" for converting map[string]interface{} to DeviceCommandMessage  return error if can't convert
func (message *DeviceCommandMessage) FromMap(val map[string]interface{}) error {
	// TODO: move this verify code to each respective sub code
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}

	if !common.StringInSlice(message.RelayID, common.PossibleRelays) {
		return errors.New("Invalid Relay ID")
	} else {
		switch message.State.Mode {
		case "manual":
			if detail, ok := message.State.Detail.(string); !ok || (detail != "on" && detail != "off") {
				return errors.New("Invalid Detail for manual")
			}
		case "auto":
			str, _ := json.Marshal(message.State.Detail)
			var thresArray device.Condition
			err := json.Unmarshal(str, &thresArray)
			if err != nil || !thresArray.Validate() {
				return errors.New("Invalid detail for auto")
			}
		case "scheduled":
			var sched device.ScheduleDetail
			str, err := json.Marshal(message.State.Detail)
			if err != nil {
				return err
			}
			err = json.Unmarshal(str, &sched)
			if err != nil {
				return errors.New("Invalid Detail - Struct" + err.Error())
			} else {
				// Check schedule
				for _, entry := range sched.Schedules {
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
			}
		default:
			return errors.New("Invalid mode")
		}
	}

	return nil
}

// FromMap is "constructor" for converting map[string]interface{} to Message  return error if can't convert
func (message *Message) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

//FromMap is "constructor" for converting map[string]interface{} to RenameDeviceMessage return error if can't convert
func (message *RenameDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

//FromMap is "constructor" for converting map[string]interface{} to TimeQuery return error if can't convert
func (message *TimeQuery) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}
