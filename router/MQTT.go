// MQTT message Handler 
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


// topic is in format cufarm1.0/xxxxxx/(status/resp/command/normal)
func idFromTopic(topic []byte) string {
	return strings.Split(string(topic), "/")[1]
}

func messageType(topic []byte) string {
	return strings.Split(string(topic), "/")[2]
}
// InitMQTT sets handler of mqtt router
func InitMQTT() {
	mqtt.SetHandler(handleMessage)
}

// InitMQTTNull sets handler to "nullHandler" which does nothing
func InitMQTTNull() {
	mqtt.SetHandler(nullHandler)
}

func nullHandler(msg *message.PublishMessage) error {
	return nil
}

func handleMessage(msg *message.PublishMessage) error {
	inMessage := []byte(string(msg.Payload()))
	deviceID := idFromTopic(msg.Topic())
	common.Println("[MQTT] <<< ", string(inMessage))

	var message model.DeviceMessage
	err := json.Unmarshal(inMessage, &message)
	common.Printf("[MQTT] <<< parsed Data=%#v\n", message)

	if err == nil {
		device, err := storage.GetDevice(deviceID)
		if err != nil {
			fmt.Println("  At handleMessage : handleMessage -> GetDevice")
			return err
		}
		// common.Printf("[MQTT] --- deviceID=[%s]\n", deviceID)
		user := storage.GetUserStateInfo(device.Owner)
		// common.Printf("[MQTT] --- owner=%s\n", device.Owner)
		switch message.Type {
		case "greeting":
			device.BroadCast()
		case "data":
			// device.UpdateState(message.Payload)
			user.ReportStatus(message.Payload, device.ID)
		}
	} else {
		common.Println("[MQTT] !!! Not a data message")
	}
	return nil
}
