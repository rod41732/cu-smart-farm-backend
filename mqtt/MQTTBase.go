// this is  scaffold implemenation for MQTT server
package mqtt

import (
	"fmt"
	"time"
	"sync"

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
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("backend-main-service"))
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
func SendMessageToDevice(version, deviceID, subTopic string, payload []byte) {
	common.Printf("[MQTT-OUT] v%s d:%s %s: %s", version, deviceID, subTopic, payload)
	publishToMQTT([]byte(fmt.Sprintf("cufarm%s/%s/%s", version, deviceID, subTopic)), payload) // TODO
}

// fix concurrent map write when connecting
var connectMux sync.Mutex

// create new connection to be used to publish
func newConnection() (*service.Client, error) {

	clnt := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte(config.MQTT["username"]))
	msg.SetPassword([]byte(config.MQTT["password"]))
	msg.SetWillQos(1)
	msg.SetVersion(3)
	msg.SetClientId([]byte("single-publish-" + common.RandomString(6)))
	msg.SetCleanSession(true)
	msg.SetKeepAlive(15)
	msg.SetWillTopic([]byte("CUSmartFarm"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	connectMux.Lock()
	err := clnt.Connect(config.MQTT["address"], msg)
	connectMux.Unlock()
	if common.PrintError(err) {
		fmt.Println("  At MQTT/connectToMQTTServer")
		return nil, err
	}
	return clnt, nil
}

func publishToMQTT(topic, payload []byte) {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(1)
	msg.SetPayload([]byte(payload))
	clnt, err := newConnection()
	for ; err != nil; clnt, err = newConnection() {
		fmt.Println("[MQTT] Can't connect to server, retrying...")
		fmt.Println(" -- At mqtt/publishToMQTT")
	}
	clnt.Publish(msg, nil)
	// TODO: if we disconnect now -> server will reject all connection from our IP
	// as will close connection to our IP
	// clnt.Disconnect()
}

func subAll() error {
	common.Println("[MQTT] ---- subscribing to all topic")
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("cufarm1.0"), 2)
	subMsg.AddTopic([]byte("cufarm1.0/+/status"), 2)
	subMsg.AddTopic([]byte("cufarm1.0/+/response"), 2)
	subMsg.AddTopic([]byte("cufarm1.0/+/greeting"), 2)
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
