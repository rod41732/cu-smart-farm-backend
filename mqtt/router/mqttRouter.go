package router

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/mqtt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	"github.com/surgemq/message"
)

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		if common.PrintError(mqtt.connectToMQTTServer()) {
			fmt.Println("[ERROR] error connecting to MQTT")
			continue
		}
		subAll()
		common.ShouldPrintDebug = true
		common.BatchWriteSize = 1
		fmt.Println("Connected.")
		time.Sleep(45 * time.Second)
		fmt.Println("Reconnecting")
	}
}
func subAll() {
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm"), 2)
	subMsg.AddTopic([]byte("CUSmartFarm/+"), 2)
	common.PrintError(mqttClient.Subscribe(subMsg, nil, handleMessage))
}

func handleMessage(msg *message.PublishMessage) error {
	message := []byte(string(msg.Payload()))
	deviceID := idFromTopic(msg.Topic())
	common.Println("[MQTT] incoming message ", string(message))

	var parsedData model.DeviceMessage
	err := json.Unmarshal(message, &parsedData)
	common.PrintError(err)
	if err == nil && parsedData.Type == "greeting" {
		return greetDevice(deviceID)
	} else if err != nil || parsedData.Type != "data" {
		common.Println("[MQTT] Error: Invalid message format or non data")
		return nil
	}

	// send data to user
	device, err := common.FindDeviceByID(deviceID)
	user := storage.GetUserStateInfo(device.Owner)
	common.Printf("[MQTT] device=%s owner=%s\n", deviceID, device.Owner)
	if user != nil {
		user.ReportStatus(parsedData)
	}
	if err != nil && err.Error() != "not found" { // ignore device not found
		common.PrintError(err)
		return err
	}
	common.Printf("[MQTT] parsed Data=%#v\n", parsedData)
	out := parsedData.ToMap()
	delete(out, "t")
	common.WriteInfluxDB("air_sensor", map[string]string{"device": deviceID}, out)

	return nil
}

// idFromTopic return <DeviceID> from CUSmartFarm/<DeviceId>_svr_recv
func idFromTopic(topic []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(topic), "CUSmartFarm/"), "_svr_recv")
}

// greetDevice : send last device state to device
func greetDevice(deviceID string) error {
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return err
	}
	device, _ := common.FindDeviceByID(deviceID)
	common.Printf("[MQTT] device id %s => %#v\n", deviceID, device)
	msg, err := json.Marshal(bson.M{
		"cmd":   "set",
		"state": toDeviceStateMap(device.RelayStates),
	})
	common.PrintError(err)
	SendMessageToDevice(deviceID, msg)
	return nil
}

// SubscribeDevice : Subscribe device when user logged in and connected to websocket
func SubscribeDevice(deviceID string) {
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm/"+deviceID+"_svr_recv"), 2)
	common.PrintError(mqttClient.Subscribe(subMsg, nil, handleMessage))
}
