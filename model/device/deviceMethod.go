package device

import (
	"encoding/json"
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
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
			"state." + relayID:                        state,
			"pastState." + relayID + "." + state.Mode: state.Detail,
		}},
	}, &temp)

	if common.PrintError(err) {
		fmt.Println("  At Device::SetRelay -- Update relay data")
		return false, "Device modify error"
	}
	device.RelayStates[relayID] = state
	device.PastStates[relayID][state.Mode] = state.Detail
	device.BroadCast()
	return true, "OK"
}

// SetName sets name (display name) of device
func (device *Device) SetName(name string) (bool, string) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At Device::SetName -- Connecting to DB")
		return false, "DB connection error"
	}
	defer mdb.Close()

	var tmp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": device.ID,
	}).Apply(mgo.Change{Update: bson.M{"$set": bson.M{"name": name}}}, &tmp)
	if common.PrintError(err) {
		fmt.Println("  At Device::SetName -- Update name")
		return false, "Update name"
	}
	device.Name = name
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
	device.sendMsg(mqttMsg)
}

// Poll : send "fetch" command to device
func (device *Device) Poll() {
	mqttMsg, _ := json.Marshal(bson.M{
		"cmd": "fetch",
	})
	device.sendMsg(mqttMsg)
}

// UpdateState update device state (from auto logic) according to data
func (dev *Device) UpdateState(data model.DeviceMessagePayload) {
	// if (rand.Float64() > 0.2) {
	// 	return
	// }

	resultState := make(map[string]RelayState)
	for _, key := range common.PossibleRelays {
		//common.Printf("[Device] ID: %s Relay:%s => Mode=%s\n", dev.ID, key, dev.RelayStates[key].Mode)
		if state := dev.RelayStates[key]; state.Mode == "auto" {
			var cond Condition
			str, _ := json.Marshal(state.Detail)
			err := json.Unmarshal(str, &cond)
			//common.PrintError(err)
			if err == nil {
				common.Println(cond.Sensor)
				var val float32
				switch cond.Sensor {
				case "soil":
					val = data.Soil
				case "temp":
					val = data.Temp
				case "humidity":
					val = data.Humidity
				default:
					continue
				}
				newState := ""
				//common.Printf("%v [%v] %v is %v\n", val, cond.Symbol, cond.Trigger, (cond.Symbol == "<") ==
				//(val < cond.Trigger))
				if (cond.Symbol == "<") == (val < cond.Trigger) {
					newState = "on"
				} else {
					newState = "off"
				}
				resultState[key] = RelayState{Mode: "manual", Detail: newState}
			}
		} else {
			resultState[key] = dev.RelayStates[key]
			common.Println("[Device] mode isn't auto")
		}
	}
	str, _ := json.Marshal(bson.M{
		"cmd":   "set",
		"state": resultState,
	})
	common.Printf("[Device] ID: %s Data = %#v State = %#v\n", dev.ID, data, dev.RelayStates)
	common.Printf("[Device] ID: %s >> sending ...\n", dev.ID)
	mqtt.SendMessageToDevice(dev.ID, str)
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
