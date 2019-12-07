package storage

import (
	"errors"
	"fmt"
	"github.com/rod41732/cu-smart-farm-backend/model"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"

	"github.com/rod41732/cu-smart-farm-backend/model/user"
)

// for marshalled db
type userData struct {
	Username string   `json:"username"`
	Devices  []string `json:"devices"`
}

var mappedUserObject = make(map[string]*user.RealUser)

// SetUserStateInfo : Map username into *user.RealUser
func SetUserStateInfo(username string, user *user.RealUser) {
	fmt.Printf("added user: %s\n", username)
	mappedUserObject[username] = user
}

func ensureUser(username string) {
	common.Printf("[Storage] create new user %s\n", username)
	var tmp user.RealUser
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Storage::ensureUser -- Connecting to DB")
		return
	}
	defer mdb.Close()
	mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": username,
	}).One(&tmp)
	tmp.Init()
	mappedUserObject[username] = &tmp
}

// GetUserStateInfo get *user.RealUser corresponding to username
func GetUserStateInfo(username string) *user.RealUser {
	_, ok := mappedUserObject[username]
	if !ok {
		ensureUser(username)
	}
	return mappedUserObject[username]
}

var Devices = make(map[string]*device.Device)

func ensureDevice(deviceID string) {
	common.Println("[Storage] make new device for", deviceID)
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At User::ensureDevice -- Connecting to DB")
		return
	}
	defer mdb.Close()
	var tmp map[string]interface{}
	err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": deviceID,
	}).One(&tmp)
	if err != nil {
		return
	}
	dev := device.Device{}
	dev.FromMap(tmp)
	Devices[deviceID] = &dev
}

// LoadDevice : loads device from Database
func LoadDevice(deviceID string) error {
	mdb, err := common.Mongo()
	if err != nil {
		return err
	}
	defer mdb.Close()
	var tmp map[string]interface{}
	var dev device.Device
	// common.Printf("[Worker-storage] before %v\n", Devices[deviceID])
	err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": deviceID,
	}).One(&tmp)
	if err != nil {
		return err
	}
	lastSensorValue := model.DeviceMessageV1_0{} 
	if (Devices[deviceID] != nil){
		lastSensorValue = Devices[deviceID].LastSensorValues
	}
	dev.FromMap(tmp)
	dev.LastSensorValues = lastSensorValue
	Devices[deviceID] = &dev
	return nil
}

// GetDevice get device object
func GetDevice(deviceID string) (dev *device.Device, err error) {
	_, ok := Devices[deviceID]
	if !ok { // then make the new device
		ensureDevice(deviceID) // load device from db if possible
	}
	res := Devices[deviceID]
	if res == nil {
		return res, errors.New("Device Not found")
	} else {
		return res, nil
	}
}

func SetDevice(deviceID string, sensorData model.DeviceMessageV1_0) {
	fmt.Printf("\nset sensor %s: %#v\n", deviceID, sensorData)
	Devices[deviceID].LastSensorValues = sensorData
	fmt.Printf("\nset sensor %s: %#v\n", deviceID, Devices[deviceID].LastSensorValues)
}
