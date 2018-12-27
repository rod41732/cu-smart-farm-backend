package user

import (
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/rod41732/cu-smart-farm-backend/common"
	mMessage "github.com/rod41732/cu-smart-farm-backend/model/message"
	"github.com/rod41732/cu-smart-farm-backend/storage"
	mgo "gopkg.in/mgo.v2"
)

// AddDevice adds device into user's device list
func (user *RealUser) AddDevice(payload map[string]interface{}) (bool, string) {
	var message mMessage.AddDeviceMessage
	if message.FromMap(payload) != nil {
		return false, "Bad Request"
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return false, "!! DB Connect error"
	}

	device, err := storage.GetDevice(message.DeviceID)
	common.Printf("[User] add device -> device=%#v\n", device)
	if err != nil {
		common.PrintError(err)
		return false, "Device not found"
	}
	if device.Owner != "" {
		common.Println("device is own")
		return false, "Device already owned "
	}
	if device.SetOwner(user.Username) {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$push": bson.M{"devices": message.DeviceID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At modifying user")
			return false, "!! user modify error"
		}
		user.devices = append(user.devices, message.DeviceID)
		return true, "OK"
	}
	return false, "!! Device modiy error"
}

// RemoveDevice removes device from user's device list
func (user *RealUser) RemoveDevice(payload map[string]interface{}) (bool, string) {
	var message mMessage.RemoveDeviceMessage
	if message.FromMap(payload) != nil {
		return false, "Bad Request"
	}

	// owner check
	if !user.ownsDevice(message.DeviceID) {
		return false, "Not your device"
	}

	// DB Operations
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	device, err := storage.GetDevice(message.DeviceID)
	if err != nil {
		return false, "Device not found"
	}
	if device.Owner != user.Username {
		return false, "Not your device"
	}
	if device.RemoveOwner() {
		var temp map[string]interface{}
		_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
			"username": user.Username,
		}).Apply(mgo.Change{
			Update: bson.M{"$pull": bson.M{"devices": message.DeviceID}},
		}, &temp)
		if common.PrintError(err) {
			fmt.Println("  At modifying user")
			return false, "!! user modify error"
		}
		common.RemoveStringFromSlice(message.DeviceID, user.devices)
		return true, "OK"
	}
	return false, "!! device modify error"
}

// PollDevice send "fetch" command to device
func (user *RealUser) PollDevice(payload map[string]interface{}) (bool, string) {

	var message mMessage.RemoveDeviceMessage
	if message.FromMap(payload) != nil {
		return false, "Bad request"
	}

	device, err := storage.GetDevice(message.DeviceID)
	if err != nil {
		return false, "Device not found"
	}
	device.Poll()
	return true, "OK"
}

// SetDevice : set relay state of device (specified via payload)
func (user *RealUser) SetDevice(payload map[string]interface{}) (bool, string) {
	var msg mMessage.DeviceCommandMessage
	if msg.FromMap(payload) != nil {
		return false, "Bad request"
	}
	if !user.ownsDevice(msg.DeviceID) {
		return false, "Not your device"
	}

	device, err := storage.GetDevice(msg.DeviceID)
	if err != nil {
		return false, "Device not found"
	}
	if device.SetRelay(msg.RelayID, msg.State) {
		return true, "OK"
	}
	return false, "Something went wrong"
}
