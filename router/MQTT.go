// MQTT message Handler
package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

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

func getVersion(topic []byte) string {
	mainTopic := strings.Split(string(topic), "/")[0] // cufarmxxxx
	return strings.TrimPrefix(mainTopic, "cufarm")
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

var persistentDevice = make(map[string]bool) // true if must repeat sending

func handleMessage(msg *message.PublishMessage) error {
	inMessage := msg.Payload()
	common.Printf("[MQTT] topic: %s <<< %s", msg.Topic(), inMessage)
	version := getVersion(msg.Topic())

	switch version {
	case "1.0":
		return handleV1Message(msg)
	default:
		common.Println("[MQTT] WARNING: unknown device message version")
		return errors.New("Unknown message version")
	}
}

func handleV1Message(msg *message.PublishMessage) error {
	inMessage := msg.Payload()
	deviceID := idFromTopic(msg.Topic())
	msgType := messageType(msg.Topic())

	payload := &model.DeviceMessageV1_0{}
	err := json.Unmarshal(inMessage, payload)
	if common.PrintError(err) {
		return err
	}

	dev, err := storage.GetDevice(deviceID)

	switch msgType {
	case "response": // device now has response
		device.UrgentFlagMux.Lock();
		device.UrgentFlag[deviceID] = false
		device.UrgentFlagMux.Unlock();
		fallthrough
	case "status": // just periodic report
		storage.SetDevice(deviceID, *payload)
		fmt.Print("Received sensor value = ", *payload)
		user := storage.GetUserStateInfo(dev.Owner)
		user.ReportStatus(payload, deviceID)
		dev.BroadCast("1.0", false);
	case "greeting": // greeting when device just connected server
		fmt.Println("Get Greeting from", deviceID)
		dev.BroadCast("1.0", true)
	}

	return nil
}
