package device

// Device : device connected ti
type Device struct {
	ID          string                `json:"id" binding:"required"`
	Secret      string                `json:"secret" binding:"required"`
	Owner       string                `json:"owner" binding:"required"`
	RelayStates map[string]RelayState `json:"state" binding:"required"`
}
