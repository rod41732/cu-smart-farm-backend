package router

import (
	"bytes"
	"fmt"
	"time"

	"../common"
	"github.com/surgemq/message"
)

func handleMessage(msg *message.PublishMessage) error {
	parsedData := common.ParseJSON(bytes.Trim(msg.Payload(), "\x00"))
	fmt.Printf("%#v\n", parsedData)
	_, ok1 := parsedData["Humidity"]
	_, ok2 := parsedData["Temp"]
	if !ok1 || !ok2 {
		fmt.Println("Error, invalid data")
		return nil
	}

	// fake device name for now
	common.WriteInfluxDB("air_sensor", map[string]string{"device": "1234"}, parsedData)
	return nil
}

func handleSubscriptionComplete(msg, ack message.Message, err error) error {
	// fmt.Printf("Subscribed: %s\nAck: %s\n", msg.Decode([]byte("utf-8")), ack.Decode([]byte("utf-8")))
	fmt.Print(msg, ack)
	common.CheckErr("OnSubComplete Handler", err)
	return nil
}

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		err := common.CheckErr("connect to MQTT", common.ConnectToMQTT())

		if err {
			fmt.Println("[ERROR] error connecting to MQTT")
			continue
		}
		subMsg := message.NewSubscribeMessage()
		subMsg.AddTopic([]byte("CUSmartFarm"), 2)
		// subMsg.SetRemainingLength()
		common.CheckErr("Subscribing", common.MqttClient.Subscribe(subMsg, handleSubscriptionComplete, handleMessage))
		fmt.Println("Connected.")
		time.Sleep(55 * time.Second)
		fmt.Println("Reconnecting")
	}
}
