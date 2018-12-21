package user

import (
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// HandleAddDevice handle add device payload, return boolean indicating success
func HandleAddDevice(payload model.AddDeviceMessage, username string) (bool, string) {
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	deviceID := payload.DeviceID
	deviceSecret := common.SHA256(payload.DeviceSecret)
	col := mdb.DB("CUSmartFarm").C("devices")

	var match bson.M
	deviceQuery := col.Find(bson.M{
		"id":     deviceID,
		"secret": deviceSecret,
	})
	deviceQuery.One(&match)

	if match == nil {
		return false, "Invalid device ID/ Secret"
	}
	if match["owner"] != nil {
		return false, "Device already owned"
	}

	var result interface{}
	appendDevice := mgo.Change{
		Update: bson.M{"$push": bson.M{"ownedDevices": deviceID}},
	}
	changeOwner := mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": username}},
	}

	col.Find(bson.M{
		"username": username,
	}).Apply(appendDevice, result)
	deviceQuery.Apply(changeOwner, result)
	return true, "OK"
}

// HandleRemoveDevice handles removal of device (and check owner before doing so)
func HandleRemoveDevice(payload model.RemoveDeviceMessage, username string) (bool, string) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	deviceID := payload.DeviceID
	col := mdb.DB("CUSmartFarm").C("devices")

	var match bson.M
	deviceQuery := col.Find(bson.M{
		"id":    deviceID,
		"owner": username,
	})
	deviceQuery.One(&match)

	if match == nil {
		return false, "Invalid DeviceID or Not Your Device"
	}

	var result interface{}
	removeDevice := mgo.Change{
		Update: bson.M{"$pull": bson.M{"ownedDevices": deviceID}},
	}
	changeOwner := mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": nil}},
	}

	col.Find(bson.M{
		"username": username,
	}).Apply(removeDevice, result)
	deviceQuery.Apply(changeOwner, result)

	return true, "OK"
}
