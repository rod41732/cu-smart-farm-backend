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
	DeviceDesc   string `json:"deviceDesc"`
}

// RenameRelayMessage is payload format for renameRelay API
type RenameRelayMessage struct {
	RelayID     string `json:"relayID"`
	Description string `json:"desc"`
}

// DeviceCommandMessage is payload format for setDevice API
type DeviceCommandMessage struct {
	RelayID string            `json:"relayID"`
	State   device.RelayState `json:"state"`
}

// EditDeviceMessage is payload format for rename device API
type EditDeviceMessage struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
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
			_map, ok := message.State.Detail.(map[string]interface{})
			if !ok {
				return errors.New("Invalid detail : isn't map")
			}
			err = sched.FromMap(_map)
			if err != nil {
				return errors.New("Invalid Detail - Struct" + err.Error())
			}
			return nil
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
func (message *EditDeviceMessage) FromMap(val map[string]interface{}) error {
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

//FromMap is "constructor" for converting map[string]interface{} to RenameRelayMessage return error if can't convert
func (message *RenameRelayMessage) FromMap(val map[string]interface{}) error {
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
