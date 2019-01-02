package storage

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"

	"github.com/rod41732/cu-smart-farm-backend/model"
)

var mappedUserObject = make(map[string]model.User)

// SetUserStateInfo : Map username into model.User
func SetUserStateInfo(username string, user model.User) {
	fmt.Printf("added user: %s\n", username)
	mappedUserObject[username] = user
}

// GetUserStateInfo get model.User corresponding to username
func GetUserStateInfo(username string) model.User {
	_, ok := mappedUserObject[username]
	fmt.Printf("[Storage]get user: %s is ok=%v\n", username, ok)
	// if !ok {
	// 	mappedUserObject[username] = &user.NullUser{}
	// }
	return mappedUserObject[username]
}

var mappedDeviceObject = make(map[string]*device.Device)

// GetDevice get device object
func GetDevice(deviceID string) (dev *device.Device, err error) {
	_, ok := mappedDeviceObject[deviceID]
	if !ok { // then make the new device
		common.Println("make new device")
		mdb, err := common.Mongo()
		if common.PrintError(err) {
			fmt.Println("  At GetDevice()")
			return &device.Device{}, err
		}
		var tmp map[string]interface{}
		mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
			"id": deviceID,
		}).One(&tmp)
		dev := device.Device{}
		dev.FromMap(tmp)
		mappedDeviceObject[deviceID] = &dev

	}
	return mappedDeviceObject[deviceID], nil
}
