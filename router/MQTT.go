package router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/rod41732/cu-smart-farm-backend/storage"
	"github.com/surgemq/message"
)

func idFromTopic(topic []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(topic), "CUSmartFarm/"), "_svr_recv")
}

// greetDevice : send last device state to device
func greetDevice(deviceID string) error {
	dev, err := storage.GetDevice(deviceID)
	if err != nil {
		return err
	}
	dev.BroadCast()
	return nil
}

// InitMQTT sets handler of mqtt router
func InitMQTT() {
	mqtt.SetHandler(handleMessage)
}

func handleMessage(msg *message.PublishMessage) error {
	topic := string(msg.Topic())
	if strings.HasSuffix(topic, "svr_out") { // skip out message
		return nil
	}

	inMessage := []byte(string(msg.Payload()))
	deviceID := idFromTopic(msg.Topic())
	common.Println("[MQTT] <<< ", string(inMessage))

	var message model.DeviceMessage
	err := json.Unmarshal(inMessage, &message)
	common.Printf("[MQTT] <<< parsed Data=%#v\n", message)

	common.PrintError(err)
	if err == nil && message.Type == "greeting" {
		return greetDevice(deviceID)
	} else if err != nil || message.Type != "data" {
		common.Println("[MQTT] !!! Not a data message")
		return nil
	}

	// send data to user
	device, err := storage.GetDevice(deviceID)
	if common.PrintError(err) {
		fmt.Println("  At handleMessage : greetDevice")
		return err
	}
	common.Printf("[MQTT] --- deviceID=[%s]\n", deviceID)
	user := storage.GetUserStateInfo(device.Owner)
	common.Printf("[MQTT] --- owner=%s\n", device.Owner)
	if user != nil {
		user.ReportStatus(message)
	}
	if err != nil && err.Error() != "not found" { // ignore device not found
		common.PrintError(err)
		return err
	}
	out := message.ToMap()
	delete(out, "t")
	common.WriteInfluxDB("air_sensor", map[string]string{"device": deviceID}, out)

	return nil
}
