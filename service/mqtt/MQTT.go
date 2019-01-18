package mqtt

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var mqttClient *service.Client

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
	err := mqttClient.Connect(config.MQTT["address"], msg)
	if common.PrintError(err) {
		fmt.Println("  At MQTT/connectToMQTTServer")
		mqttClient = nil
		return err
	}
	return nil
}

// SendMessageToDevice : Shorthand for creating message and publish
func SendMessageToDevice(deviceID string, payload []byte) {
	common.Printf("[MQTT] >>> send message: %s to %s\n", string(payload), deviceID)
	publishToMQTT([]byte("CUSmartFarm/"+deviceID+"/svr_out"), payload)
	// publishToMQTT([]byte("CUSmartFarm"), payload)
}

func publishToMQTT(topic, payload []byte) {
	fmt.Println("[MQTT]Connecting to server")
	for common.PrintError(connectToMQTTServer()) {
		fmt.Println("  At MQTT/MQTT -- Connecting to server")
		fmt.Println("[MQTT] Failed to connect to server, retrying...")
	}
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(1)
	msg.SetPayload([]byte(payload))
	mqttClient.Publish(msg, nil)
	// mqttClient.Disconnect()
}

// // MQTT : intialize MQTT Client
// func MQTT() error {
// 	for {
// 		if common.PrintError(connectToMQTTServer()) {
// 			fmt.Println("  At MQTT/MQTT -- Connecting to server")
// 			fmt.Println("[MQTT] Failed to connect to server")
// 		}
// 		publishToMQTT([]byte("CUSmartFarm"), []byte("Mother fucker"))
// 		time.Sleep(45 * time.Second)
// 		fmt.Println("[MQTT] ---- Reconnecting.")
// 	}
// }
