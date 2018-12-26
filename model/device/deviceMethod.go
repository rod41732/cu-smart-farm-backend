package device

import (
	"encoding/json"
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/mqtt"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Must check before call

// SetOwner sets device owner
func (device *Device) SetOwner(newOwner string) bool {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At device.setOwner() - db connect")
		return false
	}
	db := mdb.DB("CUSmartFarm")

	err = db.C("devices").Update(bson.M{
		"id": device.ID,
	}, mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": newOwner}},
	})
	if err != nil {
		err = db.C("users").Update(bson.M{
			"username": newOwner,
		}, mgo.Change{
			Update: bson.M{"$push": bson.M{"devices": device.ID}},
		})
	}

	if common.PrintError(err) {
		fmt.Println("  At device.setOwner) - updating")
		return false
	}
	device.Owner = newOwner
	return true
}

// RemoveOwner removes device owner
func (device *Device) RemoveOwner() bool {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At device.setOwner() - db connect")
		return false
	}
	db := mdb.DB("CUSmartFarm")

	err = db.C("devices").Update(bson.M{
		"id": device.ID,
	}, mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": nil}}, // it's set to null in DB
	})
	if err != nil {
		err = db.C("users").Update(bson.M{
			"username": "",
		}, mgo.Change{
			Update: bson.M{"$pull": bson.M{"devices": device.ID}},
		})
	}

	if common.PrintError(err) {
		fmt.Println("  At device.setOwner) - updating")
		return false
	}
	device.Owner = ""
	return true
}

// SetRelay set state of relay, and broadcast change to device
func (device *Device) SetRelay(relayID string, state RelayState) bool {
	if !state.Verify() {
		return false
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At device.setRelay() - db connect")
		return false
	}

	err = mdb.DB("CUSmartFarm").C("devices").Update(bson.M{
		"id": device.ID,
	}, mgo.Change{
		Update: bson.M{"$set": bson.M{
			"state." + relayID: state,
		}},
	})

	if common.PrintError(err) {
		fmt.Println("  At device.setOwner) - updating")
		return false
	}
	device.RelayStates[relayID] = state
	device.BroadCast()
	return true
}

// BroadCast : send current state to device via MQTT
func (device *Device) BroadCast() {
	mqttMsg, err := json.Marshal(bson.M{
		"cmd":   "set",
		"state": toDeviceStateMap(device.RelayStates),
	})
	common.Printf("[BroadCast] deviceMap is %#v\n", toDeviceStateMap(device.RelayStates))
	if common.PrintError(err) {
		fmt.Println("  at device.Broadcast()")
		return
	}
	device.sendMsg(mqttMsg)
}

// Poll : send "fetch" command to device
func (device *Device) Poll() {
	mqttMsg, _ := json.Marshal(bson.M{
		"cmd": "fetch",
	})
	device.sendMsg(mqttMsg)
}

// send message to device
func (device *Device) sendMsg(payload []byte) {
	mqtt.SendMessageToDevice(device.ID, payload)
}

// Utility function
func toDeviceStateMap(relayStateMap map[string]RelayState) map[string]RelayState {
	result := make(map[string]RelayState)
	for k, v := range relayStateMap {
		result[k] = v.ToDeviceState()
	}
	return result
}
