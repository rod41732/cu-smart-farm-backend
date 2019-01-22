package mqtt

import (
	"fmt"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var mqttClient *service.Client
var messageHandler service.OnPublishFunc

func handleSubscriptionComplete(msg, ack message.Message, err error) error {
	// fmt.Printf("Subscribed: %s\nAck: %s\n", msg.Decode([]byte("utf-8")), ack.Decode([]byte("utf-8")))
	return nil

	common.Println(msg)
	common.Println(ack)
	if common.PrintError(err) {
		fmt.Println("  At MQTT/handleSubscriptionComplete")
	}
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
	msg.SetClientId([]byte(common.RandomString(32)))
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

// SetHandler sets handler for mqtt message, must be called before connection
func SetHandler(handler service.OnPublishFunc) {
	messageHandler = handler
}

// SendMessageToDevice : Shorthand for creating message and publish
func SendMessageToDevice(deviceID string, payload []byte) {
	common.Printf("[MQTT] >>> send message: %s to %s\n", string(payload), deviceID)
	publishToMQTT([]byte("CUSmartFarm/"+deviceID+"/svr_out"), payload)
	// publishToMQTT([]byte("CUSmartFarm"), payload)
}

// create new connection to be used to publish
func newConnection() (*service.Client, error) {

	clnt := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte(config.MQTT["username"]))
	msg.SetPassword([]byte(config.MQTT["password"]))
	msg.SetWillQos(1)
	msg.SetVersion(3)
	msg.SetClientId([]byte(common.RandomString(32)))
	msg.SetCleanSession(true)
	msg.SetKeepAlive(45)
	msg.SetWillTopic([]byte("CUSmartFarm"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	common.Printf("monkaS %#v\n", msg)
	err := clnt.Connect(config.MQTT["address"], msg)
	if common.PrintError(err) {
		fmt.Println("  At MQTT/connectToMQTTServer")
		return nil, err
	}
	return clnt, nil
}

func publishToMQTT(topic, payload []byte) {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(0)
	msg.SetPayload([]byte(payload))
	clnt, err := newConnection()
	common.Printf("result = CL=%#v ER=%#v\n", clnt, err)
	for ; err != nil; clnt, err = newConnection() {
		fmt.Println("[MQTT] Can't connect to server, retrying...")
		fmt.Println(" -- At mqtt/publishToMQTT")
	}
	common.Printf("result2 = CL=%#v ER=%#v\n", clnt, err)
	clnt.Publish(msg, nil)
}

func subAll() error {
	common.Println("[MQTT] ---- subscribing to all topic")
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm"), 2)
	subMsg.AddTopic([]byte("CUSmartFarm/+/svr_recv"), 2)
	err := mqttClient.Subscribe(subMsg, handleSubscriptionComplete, messageHandler)
	return err
}

// MQTT : intialize MQTT Client
func MQTT() error {
	for {
		if common.PrintError(connectToMQTTServer()) {
			fmt.Println("  At MQTT/MQTT -- Connecting to server")
			fmt.Println("[MQTT] Failed to connect to server")
		} else {
			if err := subAll(); common.PrintError(err) {
				fmt.Println("  At MQTT/subAll()")
				common.Println("[MQTT] ---- Connection Failed.")
			} else {
				common.Println("[MQTT] ---- (Re)Connected.")
			}
		}
		time.Sleep(45 * time.Second)
		fmt.Println("[MQTT] ---- Reconnecting.")
	}
}
