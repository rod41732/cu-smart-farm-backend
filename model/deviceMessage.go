package model

// DeviceMessageData is data d ecoded from device's MQTT message
// it's version independent
type DeviceMessageData interface {
	Version() string
	ToInflux() map[string]interface{} // convert to influx measurement key-value
}
