package config

import "os"

var Mongo map[string]string = map[string]string{
 "address": "mongodb://mongodb:27017",
}

var Influx map[string]string = map[string]string{
 "address": "http://iot_service_pack:8086",
 "username": "staff",
 "password": "n0th1n9n0n53n5e",
}

var MQTT map[string]string = map[string]string{
 "address": "tcp://161.200.80.206:1883",
 "username": "staff",
 "password": "51<yk@3k2o18",
}


var CookieDomain string = "164.115.27.177"
var AutoPilotAddr string = "auto_pilot:5555"

// Init modify config based on environment variable
func Init() {
	if os.Getenv("MONGO_ADDRESS") != "" {
		Mongo["address"] = os.Getenv("MONGO_ADDRESS")
	}
	if os.Getenv("MQTT_ADDRESS") != "" {
		MQTT["address"] = os.Getenv("MQTT_ADDRESS")
	}
	if os.Getenv("MQTT_USERNAME") != "" {
		MQTT["username"] = os.Getenv("MQTT_USERNAME")
	}
	if os.Getenv("MQTT_PASSWORD") != "" {
		MQTT["password"] = os.Getenv("MQTT_PASSWORD")
	}
	if os.Getenv("INFLUX_ADDRESS") != "" {
		Influx["address"] = os.Getenv("INFLUX_ADDRESS")
	}
	if os.Getenv("INFLUX_USERNAME") != "" {
		Influx["username"] = os.Getenv("INFLUX_USERNAME")
	}
	if os.Getenv("INFLUX_PASSWORD") != "" {
		Influx["password"] = os.Getenv("INFLUX_PASSWORD")
	}
	if os.Getenv("COOKIE_DOMAIN") != "" {
		CookieDomain = os.Getenv("COOKIE_DOMAIN")
	}
}

