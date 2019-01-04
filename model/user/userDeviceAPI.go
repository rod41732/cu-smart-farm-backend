package user

import (
	"encoding/json"
	"fmt"

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
		return false, "!! DB Connect error"
	}
	defer mdb.Close()

	common.Printf("[User] add device -> device=%#v\n", device)
	if device.Owner != "" {
		common.Println("device is own")
		return false, "Device already owned "
	}
	// Check Device
	if device.SetOwner(user.Username, message.DeviceSecret) {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$push": bson.M{"devices": device.ID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At modifying user")
			return false, "!! user modify error"
		}
		user.Devices = append(user.Devices, device.ID)
		device.SetName(message.DeviceName)
		return true, "OK"
	}
	return false, "!! Device modiy error"
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
		return false, "Something went wrong"
	}

	if device.RemoveOwner() {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$pull": bson.M{"devices": device.ID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At modifying user")
			return false, "!! user modify error"
		}
		common.RemoveStringFromSlice(device.ID, user.Devices)
		return true, "OK"
	}
	return false, "!! device modify error"
}

// RenameDevice renames device
func (user *RealUser) RenameDevice(payload map[string]interface{}, device *device.Device) (bool, string) {
	// owner check
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}

	var message mMessage.RenameDeviceMessage
	err := message.FromMap(payload)
	if err != nil {
		return false, "Bad Payload"
	}

	// DB Operations
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}
	defer mdb.Close()

	if device.SetName(message.Name) {
		return true, "OK"
	}

	return false, "!! device modify error"
}

// SetDevice : set relay state of device (specified via `state`)
func (user *RealUser) SetDevice(state map[string]interface{}, device *device.Device) (bool, string) {
	var msg mMessage.DeviceCommandMessage
	if msg.FromMap(state) != nil {
		return false, "Bad request"
	}
	if !user.ownsDevice(device.ID) {
		return false, "Not your device"
	}
	if device.SetRelay(msg.RelayID, msg.State) {
		return true, "OK"
	}
	return false, "Something went wrong"
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

// QueryDeviceLog return device's log, if user owns the device, otherwise return empty array
func (user *RealUser) QueryDeviceLog(timeParams map[string]interface{}, device *device.Device) (bool, string, []client.Result) {
	if !user.ownsDevice(device.ID) {
		return false, "Not your device", make([]client.Result, 0)
	}
	var msg mMessage.TimeQuery
	if common.PrintError(msg.FromMap(timeParams)) {
		return false, "Bad request can't parse", make([]client.Result, 0)
	}
	if msg.Limit <= 0 {
		msg.Limit = 1
	} else if msg.Limit > 100 {
		msg.Limit = 100
	}
	return true, "OK", common.QueryInfluxDB(fmt.Sprintf(`SELECT *::field FROM air_sensor WHERE "device"='%s' and "time" > %v and "time" < %v ORDER BY time DESC LIMIT %d`, device.ID, msg.From.UnixNano(), msg.To.UnixNano(), msg.Limit))
}
