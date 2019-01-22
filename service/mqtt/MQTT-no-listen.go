package mqtt

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func connectToMQTTServer() (error, *service.Client) {
	mqttClient := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte(config.MQTT["username"]))
	msg.SetPassword([]byte(config.MQTT["password"]))
	msg.SetWillQos(1)
	msg.SetVersion(3)
	msg.SetClientId([]byte("backend-autopilot-" + common.RandomString(10)))
	msg.SetCleanSession(true)
	msg.SetKeepAlive(45)
	msg.SetWillTopic([]byte("CUSmartFarm"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	err := mqttClient.Connect(config.MQTT["address"], msg)
	if common.PrintError(err) {
		fmt.Println("  At MQTT/connectToMQTTServer")
		mqttClient = nil
		return err, nil
	}
	return nil, mqttClient
}

// SendMessageToDevice : Shorthand for creating message and publish
func SendMessageToDevice(deviceID string, payload []byte) {
	common.Printf("[MQTT] >>> send message: %s to %s\n", string(payload), deviceID)
	publishToMQTT([]byte("CUSmartFarm/"+deviceID+"/svr_out"), payload)
	// publishToMQTT([]byte("CUSmartFarm"), payload)
}

func publishToMQTT(topic, payload []byte) {
	fmt.Println("[MQTT]Connecting to server")
	var clnt *service.Client
	err, clnt := connectToMQTTServer()
	for ; common.PrintError(err); err, clnt = connectToMQTTServer() {
		fmt.Println("  At MQTT/MQTT -- Connecting to server")
		fmt.Println("[MQTT] Failed to connect to server, retrying...")
	}
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(1)
	msg.SetPayload([]byte(payload))
	clnt.Publish(msg, nil)
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
