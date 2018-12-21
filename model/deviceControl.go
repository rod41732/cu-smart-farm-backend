package model

type APICall struct {
	EndPoint string                 `json:"endPoint" binding:"required"` // addDevice, removeDevice, setDevice, pollDevice, listDevice...
	Payload  map[string]interface{} `json:"payload" binding:"required"`  // json data depend on command
}

// DeviceCommand :
type DeviceCommand struct {
	CmdName  string `json:"cmdName"` // such as Set, ReportStatus
	DeviceID string `json:"deviceID"`
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
