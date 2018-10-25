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
	fmt.Println(parsedData)
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
		common.CheckErr("connect to MQTT", common.ConnectToMQTT())

		subMsg := message.NewSubscribeMessage()
		subMsg.AddTopic([]byte("CUSmartFarm"), 1)
		common.CheckErr("Subscribing", common.MqttClient.Subscribe(subMsg, handleSubscriptionComplete, handleMessage))

		fmt.Print("Connected.")
		time.Sleep(10 * time.Second)
		fmt.Print("Reconnecting")
	}
}
