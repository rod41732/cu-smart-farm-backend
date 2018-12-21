package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../config"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"gopkg.in/mgo.v2"
)

const (
	privKeyPath = "key.rsa"
	pubKeyPath  = "key.rsa.pub"
)

var (
	// SignKey = private key
	SignKey []byte

	// VerifyKey = public key
	VerifyKey []byte
)

// MqttClient : this is MQTT client that listen to server
var MqttClient *service.Client

// BatchWriteSize : How many points to write at once (set to 1 isn't a problem)
var BatchWriteSize = 3

// ShouldPrintDebug this flag control whether we should print debug
var ShouldPrintDebug = false

// InitializeKeyPair initializes public/private key pair
func InitializeKeyPair() {
	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}
}

// SHA256 : sha256 encrypt helper
func SHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// PrintError : return true and print if error
func PrintError(err error) bool {
	if err != nil {
		fmt.Printf("[Error] %s\n", err)
		return true
	}
	return false
}

// Resp : response with that status and return true if error
func Resp(statusCode int, err error, c *gin.Context) bool {
	if err != nil {
		if !ShouldPrintDebug {
			c.JSON(statusCode, "something went wrong")
			return true
		}
		c.JSON(statusCode, err.Error())
		return true
	}
	return false
}

// Mongo returns a session
func Mongo() (*mgo.Session, error) {
	return mgo.Dial(config.Mongo["address"])
}

// ConnectToMQTT : connects to mqtt server and return error if error
func ConnectToMQTT() error {
	if MqttClient != nil {
		MqttClient.Disconnect()
	}
	MqttClient = &service.Client{}

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
	PrintError(MqttClient.Connect(config.MQTT["address"], msg))
	// msg.SetCleanSession(true)
	return nil
}

// PublishToMQTT : Shorthand for creating message and publish
func PublishToMQTT(topic, payload []byte) {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(0)
	msg.SetPayload([]byte(payload))
	MqttClient.Publish(msg, nil)
}

// ParseJSON : parse byte to json (gin.H)
func ParseJSON(payload []byte) map[string]interface{} {
	var jsonData map[string]interface{}
	// fmt.Println("In === ", json.Unmarshal(payload, jsonData))
	json.Unmarshal(payload, &jsonData)
	return jsonData
}

// ConnectToInfluxDB : connect to influx DB and return client
func ConnectToInfluxDB() (client.Client, error) {
	influxConn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Influx["address"],
		Username: config.Influx["username"],
		Password: config.Influx["password"],
	})
	PrintError(err)
	return influxConn, err
}

// QueryInfluxDB : runs query in influxDB
func QueryInfluxDB(query string) []client.Result {
	clnt, err := ConnectToInfluxDB()

	PrintError(err)
	if err == nil {
		resp, err := clnt.Query(client.Query{
			Command:  query,
			Database: "CUSmartFarm",
		})
		PrintError(err)
		if err == nil {
			fmt.Printf("Query Success: %v \n", resp)
		}
		return resp.Results
	}
	return []client.Result{}
}

var deferredPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
	Database:  "CUSmartFarm",
	Precision: "ms",
})

// WriteInfluxDB : (deferred) Write a data point in to influxDB
func WriteInfluxDB(measurement string, tags map[string]string, fields map[string]interface{}) error {

	clnt, err := ConnectToInfluxDB()
	defer clnt.Close()
	PrintError(err)
	if err == nil {

		point, err := client.NewPoint("air_sensor", tags, fields, time.Now())
		if PrintError(err) {
			return err
		}
		deferredPoints.AddPoint(point)
		if ln := len(deferredPoints.Points()); ln < 3 {
			Printf("write deferred %d/3 points\n", ln)
		}
		if len(deferredPoints.Points()) >= 3 {
			err = clnt.Write(deferredPoints)
			if !PrintError(err) {
				Println("DB Write Succeeded", err)
			}
			// create new batch to remove all points
			deferredPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "CUSmartFarm",
				Precision: "ms",
			})
			if PrintError(err) {
				return nil
			}
		}
		return nil
	}
	return nil
}

// Print : this is literally fmt.Print but only print when ShouldPrintDebug flag is true
func Print(a ...interface{}) {
	if ShouldPrintDebug {
		fmt.Print(a...)
	}
}

// Println : this is literally fmt.Println but only print when ShouldPrintDebug flag is true
func Println(a ...interface{}) {
	if ShouldPrintDebug {
		fmt.Println(a...)
	}
}

// Printf : this is literally fmt.Printf but only print when ShouldPrintDebug flag is true
func Printf(format string, a ...interface{}) {
	if ShouldPrintDebug {
		fmt.Printf(format, a...)
	}
}

var wsClients = make(map[string]*websocket.Conn)
var wsDevices = make(map[string]*websocket.Conn)

// AddClientConn : add Client to broadcasted
func AddClientConn(deviceId string, conn *websocket.Conn) {
	wsClients[deviceId] = conn
}

// RemoveClientConn : remove Client to be broadcasted
func RemoveClientConn(deviceId string, conn *websocket.Conn) {
	wsClients[deviceId] = nil
}

// AddDeviceConn : add Device to be sent to
func AddDeviceConn(deviceId string, conn *websocket.Conn) {
	wsDevices[deviceId] = conn
}

// RemoveDeviceConn : remove Device to be sent to
func RemoveDeviceConn(deviceId string, conn *websocket.Conn) {
	wsDevices[deviceId] = nil
}

func TellDevice(deviceId string) bool {
	if conn, ok := wsDevices[deviceId]; ok {
		// TODO: chage type
		conn.WriteMessage(1, []byte(`{"t": "cmd", "cmd": "fetch"}`))
	}
}
