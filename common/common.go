package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"gopkg.in/mgo.v2"
)

// MqttClient : this is MQTT client that listen to server
var MqttClient *service.Client

// BatchWriteSize : How many points to write at once (set to 1 isn't a problem)
var BatchWriteSize = 3
var ShouldPrintDebug = true

// CheckErr : return true and print if error
func CheckErr(source string, err error) bool {
	if err != nil {
		fmt.Printf("\n[ERR] IN %s || %s\n", source, err)
		return true
	}
	return false
}

func ErrResp(source string, err error, c *gin.Context) {
	if CheckErr(source, err) {
		if !ShouldPrintDebug {
			c.JSON(500, "something went wrong")
		} else {
			c.JSON(500, "At"+source+":"+err.Error())
		}
	}
}

// Mongo returns a session
func Mongo() (*mgo.Session, error) {
	return mgo.Dial("mongodb://127.0.0.1:27017")
}

// ConnectToMQTT : connects to mqtt server and return error if error
func ConnectToMQTT() error {
	if MqttClient != nil {
		MqttClient.Disconnect()
	}
	MqttClient = &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte("admin"))
	msg.SetPassword([]byte("iyddyoot"))
	msg.SetWillQos(2)
	msg.SetVersion(3)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("backend"))
	msg.SetKeepAlive(45)
	msg.SetWillTopic([]byte("CUSmartFarm"))
	msg.SetWillMessage([]byte("backend: connecting.."))
	CheckErr("Connecting to MQTT", MqttClient.Connect("tcp://164.115.27.177:1883", msg))
	// msg.SetCleanSession(true)
	return nil
}

// PublishToMQTT : Shorthand for creating message and publish
func PublishToMQTT(topic, payload []byte) {
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte(topic))
	msg.SetQoS(0)
	text, _ := json.Marshal(payload)
	msg.SetPayload([]byte(text))
	MqttClient.Publish(msg, nil)
}

// ParseJSON : parse byte to json (gin.H)
func ParseJSON(payload []byte) map[string]interface{} {
	var jsonData map[string]interface{}
	// fmt.Println("In === ", json.Unmarshal(payload, jsonData))
	CheckErr("Parsing JSON", json.Unmarshal(payload, &jsonData))
	return jsonData
}

// ConnectToInfluxDB : connect to influx DB and return client
func ConnectToInfluxDB() (client.Client, error) {
	influxConn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "admin",
		Password: "4fs,mg-0zv",
	})
	CheckErr("Influx Connection", err)
	return influxConn, err
}

// QueryInfluxDB : runs query in influxDB
func QueryInfluxDB(query string) []client.Result {
	clnt, err := ConnectToInfluxDB()

	CheckErr("connect to query influx", err)
	if err == nil {
		resp, err := clnt.Query(client.Query{
			Command:  query,
			Database: "CUSmartFarm",
		})
		CheckErr("querying influx", err)
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
	CheckErr("connect to query influx", err)
	if err == nil {

		point, err := client.NewPoint("air_sensor", tags, fields, time.Now())
		if CheckErr("create new point", err) {
			return err
		}
		deferredPoints.AddPoint(point)
		if ln := len(deferredPoints.Points()); ln < 3 {
			Printf("write deferred %d/3 points\n", ln)
		}
		if len(deferredPoints.Points()) >= 3 {
			err = clnt.Write(deferredPoints)
			if !CheckErr("querying influx", err) {
				Println("DB Write Succeeded", err)
			}
			// create new batch to remove all points
			deferredPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "CUSmartFarm",
				Precision: "ms",
			})
			if CheckErr("Recreating batch points", err) {
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
