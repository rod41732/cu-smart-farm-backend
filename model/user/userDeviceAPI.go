package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/common"
	mMessage "github.com/rod41732/cu-smart-farm-backend/model/message"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddDevice adds device into user's device list
func (user *RealUser) AddDevice(param map[string]interface{}, device *device.Device) (bool, string) {
	var message mMessage.AddDeviceMessage
	if message.FromMap(param) != nil {
		return false, "Bad Request"
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At User::AddDevice -- Connecting to DB")
		return false, "Can't connect to DB"
	}
	defer mdb.Close()

	common.Printf("[User] add device -> device=%#v\n", device)
	if device.Owner != "" {
		common.Println("device is own")
		return false, "Device already owned"
	}
	// Check Device
	if ok, errmsg := device.SetOwner(user.Username, message.DeviceSecret); ok {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$push": bson.M{"devices": device.ID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At User::AddDevice -- Updating Device list")
			return false, "User modify Error"
		}
		user.Devices = append(user.Devices, device.ID)
		device.SetInfo(message.DeviceName, message.DeviceDesc)
		return true, "OK"
	} else {
		return false, "Device modify error " + errmsg
	}
}

// RemoveDevice removes device from user's device list
func (user *RealUser) RemoveDevice(device *device.Device) (bool, string) {
	// owner check
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}

	// DB Operations
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		fmt.Println("  At User::RemoveDevice -- Connecting to DB")
		return false, "Can't connect to DB"
	}

	if ok, errmsg := device.RemoveOwner(); ok {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$pull": bson.M{"devices": device.ID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At User::RemoveDevice -- Updating Device list")
			return false, "User modify Error"
		}
		common.RemoveStringFromSlice(device.ID, &user.Devices)
		return true, "OK"
	} else {
		return false, "Device modify error " + errmsg
	}
}

// EditDevice renames device and change description
func (user *RealUser) EditDevice(payload map[string]interface{}, device *device.Device) (bool, string) {
	// owner check
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}

	var message mMessage.EditDeviceMessage
	err := message.FromMap(payload)
	if err != nil {
		return false, "Bad Payload"
	}

	if ok, errmsg := device.SetInfo(message.Name, message.Description); ok {
		return true, "OK"
	} else {
		return false, "Device modify error" + errmsg
	}

}

// SetDeviceRelay : set relay state of device (specified via `state`)
func (user *RealUser) SetDeviceRelay(state map[string]interface{}, device *device.Device) (bool, string) {
	var msg mMessage.DeviceCommandMessage
	if err := msg.FromMap(state); err != nil {
		// common.PrintError(msg.FromMap(state))
		return false, "Bad payload " + err.Error()
	}
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}
	fmt.Println("[User] state = ", msg.State)
	if ok, errmsg := device.SetRelay(msg.RelayID, msg.State); ok {
		return true, "OK"
	} else {
		return false, "Device modify error " + errmsg
	}
}

// SetDeviceRelayName : set relay's name
func (user *RealUser) SetDeviceRelayName(payload map[string]interface{}, device *device.Device) (bool, string) {
	var msg mMessage.RenameRelayMessage
	if err := msg.FromMap(payload); err != nil {
		return false, "Bad payload " + err.Error()
	}
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}
	if ok, errmsg := device.SetRelayName(msg.RelayID, msg.Description); ok {
		return true, "OK"
	} else {
		return false, "Device modify error " + errmsg
	}
}

// GetDeviceInfo returns devices state, if user owns the device, otherwise return nil
func (user *RealUser) GetDeviceInfo(device *device.Device) (bool, string, map[string]interface{}) {
	// owner check
	if !user.ownsDevice(device.ID) {
		return false, "Not your device", nil
	}

	var result map[string]interface{}
	str, _ := json.Marshal(device)
	json.Unmarshal(str, &result)

	return true, "OK", result
}

// GetRelayNames return all relay names as []string
func (user *RealUser) GetRelayNames(device *device.Device) (bool, string, []string) {
	if !user.ownsDevice(device.ID) {
		return false, "Not your device", make([]string, 0)
	}
	result := make([]string, 5)
	for idx, key := range common.PossibleRelays {
		result[idx] = device.RelayStates[key].Description
	}
	return true, "OK", result
}

// QueryDeviceLog return device's log, if user owns the device, otherwise return empty array
func (user *RealUser) QueryDeviceLog(timeParams map[string]interface{}, device *device.Device) (bool, string, []client.Result) {
	if !user.ownsDevice(device.ID) {
		return false, "Not your device", nil
	}
	var msg mMessage.TimeQuery
	if err := msg.FromMap(timeParams); err != nil {
		return false, "Bad Payload " + err.Error(), nil
	}
	if msg.Limit <= 0 {
		msg.Limit = 10
	} /*else if msg.Limit > 100 {
		msg.Limit = 100
	}*/
	var res []client.Result
	if msg.From.IsZero() || msg.To.IsZero() { // when user just want ot get latest
		res = common.QueryInfluxDB(fmt.Sprintf(`SELECT *::field FROM deviceData WHERE "device"='%s' ORDER BY time DESC LIMIT %d`, device.ID, msg.Limit))
	} else {
		res = common.QueryInfluxDB(fmt.Sprintf(`SELECT *::field FROM deviceData WHERE "device"='%s' and "time" > %v and "time" < %v ORDER BY time DESC LIMIT %d`, device.ID, msg.From.UnixNano(), msg.To.UnixNano(), msg.Limit))
	}
	if common.HaveSeries(res) {
		for _, row := range res[0].Series[0].Values {
			timestamp, _ := time.Parse(time.RFC3339, row[0].(string))
			row[0] = timestamp.Unix()
		}
	}
	return true, "OK", res
}
