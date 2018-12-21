package router

import (
	"bytes"
	"fmt"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/surgemq/message"
)

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

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		if common.PrintError(common.ConnectToMQTT()) {
			fmt.Println("[ERROR] error connecting to MQTT")
			continue
		}
		common.ShouldPrintDebug = true
		subMsg := message.NewSubscribeMessage()
		subMsg.AddTopic([]byte("CUSmartFarm"), 2)
		// subMsg.SetRemainingLength()
		common.BatchWriteSize = 1
		common.PrintError(common.MqttClient.Subscribe(subMsg, handleSubscriptionComplete, handleMessage))
		fmt.Println("Connected.")
		time.Sleep(45 * time.Second)
		fmt.Println("Reconnecting")
	}
}
