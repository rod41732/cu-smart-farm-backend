package common

import (
	"fmt"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var MqttClient *service.Client

func CheckErr(source string, err error) bool {
	if err != nil {
		fmt.Printf("[ERR] IN %s || %s\n", source, err)
		return true
	}
	return false
}

// ConnecttToMQTT : connects to mqtt server and return error if error
func ConnectToMQTT() error {
	if MqttClient != nil {
		MqttClient.Disconnect()
	}
	MqttClient = &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername("admin")
	msg.SetPassword("iyddyoot")
	msg.SetWillQos(2)
	msg.SetKeepAlive(120)
	msg.SetVersion(3)
	msg.SetClientId([]byte("smart-farm-backend"))
	msg.SetWillTopic([]byte("backend-service"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	// msg.SetCleanSession(true)
	return MqttClient.Connect("tcp://164.115.27.177:1883", msg)
}
