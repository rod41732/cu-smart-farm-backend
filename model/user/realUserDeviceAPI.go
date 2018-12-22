package user

import (
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/rod41732/cu-smart-farm-backend/common"
	mMessage "github.com/rod41732/cu-smart-farm-backend/model/message"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	mgo "gopkg.in/mgo.v2"
)

// AddDevice adds device into user's device list
func (user *RealUser) AddDevice(payload map[string]interface{}) (bool, string) {
	var message mMessage.AddDeviceMessage
	if message.FromMap(payload) != nil {
		return false, "Bad Request"
	}
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	deviceID := message.DeviceID
	deviceSecret := common.SHA256(message.DeviceSecret)
	db := mdb.DB("CUSmartFarm")

	// Find device and update
	var match bson.M
	deviceCondition := bson.M{"id": deviceID, "secret": deviceSecret}
	db.C("devices").Find(deviceCondition).One(&match)

	if match == nil {
		return false, "Invalid device ID/ Secret"
	}
	if match["owner"] != nil {
		return false, "Device already owned"
	}

	appendDevice := mgo.Change{
		Update: bson.M{"$push": bson.M{"devices": deviceID}},
	}
	changeOwner := mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": user.Username}},
	}

	var temp map[string]interface{}
	_, err1 := db.C("users").Find(bson.M{"username": user.Username}).Apply(appendDevice, &temp)
	_, err2 := db.C("devices").Find(deviceCondition).Apply(changeOwner, &temp)
	if !common.PrintError(err1) && !common.PrintError(err2) {
		user.devices = append(user.devices, deviceID)
		return true, "OK"
	}
	return false, "Something went wrong"
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

	deviceID := message.DeviceID
	db := mdb.DB("CUSmartFarm")

	// no need to check owner as it's already checked
	deviceCondition := bson.M{"id": deviceID}
	if cnt, _ := db.C("devices").Find(deviceCondition).Count(); cnt == 0 {
		return false, "Device not found"
	}

	removeDevice := mgo.Change{
		Update: bson.M{"$pull": bson.M{"devices": deviceID}},
	}
	changeOwner := mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": nil}},
	}

	var temp map[string]interface{}
	_, err1 := db.C("users").Find(bson.M{"username": user.Username}).Apply(removeDevice, &temp)
	_, err2 := db.C("devices").Find(deviceCondition).Apply(changeOwner, &temp)
	if !common.PrintError(err1) && !common.PrintError(err2) {
		common.RemoveStringFromSlice(deviceID, user.devices)
		return true, "OK"
	}
	return false, "Something went wrong"
}

// PollDevice send "fetch" command to device
func (user *RealUser) PollDevice(payload map[string]interface{}) (bool, string) {

	var message mMessage.RemoveDeviceMessage
	if message.FromMap(payload) != nil {
		return false, "Bad request"
	}
	deviceID := message.DeviceID
	if !user.ownsDevice(deviceID) {
		return false, "Not your device"
	}
	mqttMessage, _ := json.Marshal(bson.M{
		"t":   "cmd",
		"cmd": "fetch",
	})
	mqtt.SendMessageToDevice(deviceID, mqttMessage)
	return true, "OK"
}

// Command set relay state of device (specified via payload)
func (user *RealUser) Command(deviceID string, relay string, workmode string, payload interface{}) (bool, string) {
	if !user.ownsDevice(deviceID) {
		return false, "Not your device"
	}

	return true, "OK"
}
