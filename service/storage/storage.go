package storage

import (
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model/device"
	"gopkg.in/mgo.v2/bson"
)

// Devices stores all devices in databases
var Devices = make(map[string]*device.Device)

// LoadDevice : loads device from Database
func LoadDevice(deviceID string) error {
	mdb, err := common.Mongo()
	if err != nil {
		return err
	}
	var tmp map[string]interface{}
	var dev device.Device
	// common.Printf("[Worker-storage] before %v\n", Devices[deviceID])
	err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": deviceID,
	}).One(&tmp)
	if err != nil {
		return err
	}
	dev.FromMap(tmp)
	Devices[deviceID] = &dev
	// common.Printf("[Worker-storage] after %v\n", Devices[deviceID])
	return nil
}

// GetDevice : return device from map (cache), create if not exist
func GetDevice(deviceID string) (*device.Device, error) {
	_, ok := Devices[deviceID]
	if !ok {
		err := LoadDevice(deviceID)
		if err != nil {
			return nil, err
		}
	}
	return Devices[deviceID], nil
}
