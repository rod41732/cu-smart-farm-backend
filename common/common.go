package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/rod41732/cu-smart-farm-backend/config"
	mgo "gopkg.in/mgo.v2"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// WsCommand : Message format for websocket message
type WsCommand struct {
	Endpoint string
	Payload  interface{}
}

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

// BatchWriteSize : How many points to write at once (set to 1 isn't a problem)
var BatchWriteSize = 3

// ShouldPrintDebug this flag control whether we should print debug
var ShouldPrintDebug = false

// Secure this flag control whether we check token of request
var Secure = true

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

// Mongo returns a session
func Mongo() (*mgo.Session, error) {
	return mgo.Dial(config.Mongo["address"])
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
			Database: "SkyhawkPhase1",
		})
		PrintError(err)
		if err == nil {
			Printf("[Influx] Query Success. Result = %#v \n", resp)
		}
		return resp.Results
	}
	return make([]client.Result, 0)
}

var deferredPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
	Database:  "SkyhawkPhase1",
	Precision: "ms",
})

// WriteInfluxDB : (deferred) Write a data point in to influxDB
func WriteInfluxDB(measurement string, tags map[string]string, fields map[string]interface{}) error {

	clnt, err := ConnectToInfluxDB()
	defer clnt.Close()
	PrintError(err)
	if err == nil {

		point, err := client.NewPoint("deviceData", tags, fields, time.Now())
		if PrintError(err) {
			return err
		}
		deferredPoints.AddPoint(point)
		if ln := len(deferredPoints.Points()); ln < BatchWriteSize {
			Printf("[Influx] write deferred %d/%d points\n", ln, BatchWriteSize)
		}
		if len(deferredPoints.Points()) >= BatchWriteSize {
			err = clnt.Write(deferredPoints)
			if !PrintError(err) {
				Println("[Influx] DB Write Succeeded", err)
			}
			// create new batch to remove all points
			deferredPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "SkyhawkPhase1",
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

// StringInSlice check whether string is in slice
func StringInSlice(str string, slice []string) bool {
	for _, x := range slice {
		if x == str {
			return true
		}
	}
	return false
}

// RemoveStringFromSlice removes string from slice
func RemoveStringFromSlice(str string, slice *[]string) {
	for idx, x := range *slice {
		if x == str {
			(*slice)[idx] = (*slice)[len(*slice)-1]
			*slice = (*slice)[:len(*slice)-1]
			break
		}
	}
}

// RandomString : helper function for random string with custom length and charset
func RandomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
