package modelmessage

type APICall struct {
	EndPoint string                 `json:"endPoint" binding:"required"` // addDevice, removeDevice, setDevice, pollDevice, listDevice...
	Token    string                 `json:"token" binding:"required"`
	Payload  map[string]interface{} `json:"payload" binding:"required"` // json data depend on command
}

type RelayMode struct { // use when set relay mode
	Mode   string      `json:"mode"`
	Detail interface{} `json:"detail"`
}

type AddDeviceMessage struct {
	DeviceID     string `json:"deviceID" binding:"required"`
	DeviceSecret string `json:"deviceSecret" binding:"required"`
}

type RemoveDeviceMessage struct {
	DeviceID string `json:"deviceID" binding:"required"`
}
