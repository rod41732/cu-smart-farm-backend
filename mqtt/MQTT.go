package mqtt

import (
	"bytes"
	"fmt"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var mqttClient *service.Client

func handleMessage(msg *message.PublishMessage) error {
	common.Println("[Debug] incoming message ", string(msg.Payload()))
	parsedData := common.ParseJSON(bytes.Trim(msg.Payload(), "\x00"))
	if parsedData["t"] != "data" {
		return nil
	}
	parsedData = parsedData["data"].(map[string]interface{})
	_, ok1 := parsedData["Humidity"]
	_, ok2 := parsedData["Temp"]
	_, ok3 := parsedData["Soil"]
	if !ok1 || !ok2 || !ok3 {
		fmt.Println("Error, invalid data")
		return nil
	}

	// fake device name for now
	common.WriteInfluxDB("air_sensor", map[string]string{"device": "1234"}, parsedData)
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
	msg.SetWillQos(2)
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
	publishToMQTT([]byte("CUSmartFarm/"+deviceID+"_svr_out"), payload)
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

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		if common.PrintError(connectToMQTTServer()) {
			fmt.Println("[ERROR] error connecting to MQTT")
			continue
		}
		common.ShouldPrintDebug = true
		common.BatchWriteSize = 1
		fmt.Println("Connected.")
		time.Sleep(45 * time.Second)
		fmt.Println("Reconnecting")
	}
}
