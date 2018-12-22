package mqtt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/model"

	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var mqttClient *service.Client

func handleMessage(msg *message.PublishMessage) error {
	message := []byte(string(msg.Payload()))
	common.Println("[MQTT] incoming message ", string(message))
	var parsedData model.DeviceMessage
	err := json.Unmarshal(message, &parsedData)
	if err != nil || parsedData.Type != "data" {
		common.Println("[MQTT] Error: Invalid message format or non data")
		return nil
	}
	// send data to user
	deviceID := strings.TrimSuffix(strings.TrimPrefix(string(msg.Topic()), "CUSmartFarm/"), "_svr_recv")
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
	common.WriteInfluxDB("air_sensor", map[string]string{"device": deviceID}, parsedData.ToMap())

	return nil
}

func handleSubscriptionComplete(msg, ack message.Message, err error) error {
	// fmt.Printf("Subscribed: %s\nAck: %s\n", msg.Decode([]byte("utf-8")), ack.Decode([]byte("utf-8")))
	fmt.Println(msg)
	fmt.Println(ack)
	common.PrintError(err)
	return nil
}

func connectToMQTTServer() error {
	if mqttClient != nil {
		mqttClient.Disconnect()
	}
	mqttClient = &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte(config.MQTT["username"]))
	msg.SetPassword([]byte(config.MQTT["password"]))
	msg.SetWillQos(1)
	msg.SetVersion(3)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("backend"))
	msg.SetKeepAlive(45)
	msg.SetWillTopic([]byte("CUSmartFarm"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	common.PrintError(mqttClient.Connect(config.MQTT["address"], msg))
	// msg.SetCleanSession(true)
	return nil
}

// SendMessageToDevice : Shorthand for creating message and publish
func SendMessageToDevice(deviceID string, payload []byte) {
	common.Printf("send message: %s to %s\n", string(payload), deviceID)
	publishToMQTT([]byte("CUSmartFarm/"+deviceID+"_svr_out"), payload)
	// publishToMQTT([]byte("CUSmartFarm"), payload)
}

func publishToMQTT(topic, payload []byte) {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(0)
	msg.SetPayload([]byte(payload))
	mqttClient.Publish(msg, nil)
}

// SubscribeDevice : Subscribe device when user logged in and connected to websocket
func SubscribeDevice(deviceID string) {
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm/"+deviceID+"_svr_recv"), 2)
	common.PrintError(mqttClient.Subscribe(subMsg, handleSubscriptionComplete, handleMessage))
}

func subAll() {
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm"), 2)
	subMsg.AddTopic([]byte("CUSmartFarm/+"), 2)
	common.PrintError(mqttClient.Subscribe(subMsg, handleSubscriptionComplete, handleMessage))
}

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		if common.PrintError(connectToMQTTServer()) {
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
