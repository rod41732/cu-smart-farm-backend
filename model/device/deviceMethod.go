package device

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/rod41732/cu-smart-farm-backend/config"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Must check before call

// SetOwner sets device owner, called when add device
func (device *Device) SetOwner(newOwner string, secret string) (bool, string) {
	if secret != device.Secret {
		return false, "Incorrect secret"
	}
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::SetOwner -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()
	db := mdb.DB("CUSmartFarm")

	var temp map[string]interface{}
	_, err = db.C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": newOwner}},
	}, &temp)
	if common.PrintError(err) {
		fmt.Println("  At Device::SetOwner -- Updating Owner")
		return false, "Updating Owner"
	}
	device.Owner = newOwner
	return true, "OK"
}

// RemoveOwner removes device owner
func (device *Device) RemoveOwner() (bool, string) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::RemoveOwner -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()
	db := mdb.DB("CUSmartFarm")

	var temp map[string]interface{}
	_, err = db.C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": nil}}, // it's set to null in DB
	}, &temp)
	if common.PrintError(err) {
		fmt.Println("  At Device::SetOwner -- Updating Owner")
		return false, "Updating owner"
	}
	device.Owner = ""
	return true, "OK"
}

// SetRelay set state of relay, and broadcast change to device
func (device *Device) SetRelay(relayID string, state RelayState) (bool, string) {
	if !state.Verify() {
		return false, "Invaid relay data"
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::SetRelay -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()

	var temp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{
			"state." + relayID + ".mode":              state.Mode,
			"state." + relayID + ".detail":            state.Detail,
			"pastState." + relayID + "." + state.Mode: state.Detail,
		}},
	}, &temp)

	if common.PrintError(err) {
		fmt.Println("  At Device::SetRelay -- Update relay data")
		return false, "Device modify error"
	}
	device.RelayStates[relayID] = RelayState{
		Mode:        state.Mode,
		Detail:      state.Detail,
		Description: device.RelayStates[relayID].Description,
	}

	oldState := device.RelayStates[relayID]
	oldState.Mode = state.Mode
	oldState.Detail = state.Detail
	device.RelayStates[relayID] = oldState
	device.PastStates[relayID][state.Mode] = state.Detail
	device.BroadCast()
	// Trigger reload
	clnt, err := rpc.DialHTTP("tcp", config.AutoPilotAddr)
	reply := new(string)
	clnt.Call("Trigger.Update", device.ID, reply)
	common.Println("[Caller] reply = ", *reply)
	return true, "OK"
}

func (device *Device) SetRelayName(relayID string, name string) (bool, string) {
	if !common.StringInSlice(relayID, common.PossibleRelays) {
		return false, "Invalid Relay ID"
	}
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::SetRelayName -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()
	var tmp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{
			"state." + relayID + ".desc": name,
		}},
	}, &tmp)
	if common.PrintError(err) {
		fmt.Println("  At Device::SetRelayName -- Update relay desc")
		return false, "Device modify error"
	}
	// cannot assign directly to struct field in map so I do this
	cpy := device.RelayStates[relayID]
	cpy.Description = name
	device.RelayStates[relayID] = cpy
	return true, "OK"
}

// SetInfo sets name and description of device
func (device *Device) SetInfo(name string, desc string) (bool, string) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::SetInfo -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()

	var tmp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{Update: bson.M{"$set": bson.M{"name": name, "desc": desc}}}, &tmp)
	if common.PrintError(err) {
		fmt.Println("  At Device::SetInfo -- Update name and desc")
		return false, "Update name"
	}
	device.Name = name
	device.Description = desc
	return true, "OK"
}

// BroadCast : send current state to device via MQTT
func (device *Device) BroadCast() {
	mqttMsg, err := json.Marshal(bson.M{
		"cmd":   "set",
		"state": toDeviceStateMap(device.RelayStates),
	})
	common.Printf("[BroadCast] deviceMap is %#v\n", toDeviceStateMap(device.RelayStates))
	if common.PrintError(err) {
		fmt.Println("  At Device::BroadCast")
		return
	}
	device.SendMsg(mqttMsg)
}

// Poll : send "fetch" command to device
func (device *Device) Poll() {
	mqttMsg, _ := json.Marshal(bson.M{
		"cmd": "fetch",
	})
	device.SendMsg(mqttMsg)
}

// send message to device
func (device *Device) SendMsg(payload []byte) {
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
