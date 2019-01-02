package user

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/common"
	mMessage "github.com/rod41732/cu-smart-farm-backend/model/message"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddDevice adds device into user's device list
func (user *RealUser) AddDevice(secret map[string]interface{}, device *device.Device) (bool, string) {
	var message mMessage.AddDeviceMessage
	if message.FromMap(secret) != nil {
		return false, "Bad Request"
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return false, "!! DB Connect error"
	}

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
		user.devices = append(user.devices, device.ID)
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

	if device.Owner != user.Username {
		return false, "Not your device"
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
		common.RemoveStringFromSlice(device.ID, user.devices)
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
