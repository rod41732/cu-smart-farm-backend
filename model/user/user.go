package user

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
)

// RealUser represent client connected via WebSocket
type RealUser struct {
	currentToken string
	conn         *websocket.Conn

	Devices    []string `json:"devices"`
	Username   string   `json:"username"`
	Province   string   `json:"province"`
	Email      string   `json:"email"`
	NationalID string   `json:"nationalID"`
	Address    string   `json:"address"`
}

//Init  initializes user token and connection
func (user *RealUser) Init() {
	user.RegenerateToken()
	user.conn = nil
	if user.Devices == nil {
		user.Devices = make([]string, 0)
	}
}

//SetConn is setter for Conn
func (user *RealUser) SetConn(conn *websocket.Conn) {
	user.conn = conn
}

// ReportStatus sends MQTT data to user via WebSocket then insert into InfluxDB
func (user *RealUser) ReportStatus(payload model.DeviceMessagePayload, deviceID string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"t":       "report",
		"d":       deviceID,
		"payload": payload,
	})
	if user.conn != nil {
		user.conn.WriteMessage(1, resp) // 1 is text message
	} else {
		common.Println("[User] null user => ", resp)
	}
	common.WriteInfluxDB("cu_smartfarm_sensor_log", map[string]string{"device": deviceID}, payload.ToMap())
}

// RegenerateToken : Regenerate user websocket authorization token
func (user *RealUser) RegenerateToken() string {
	user.currentToken = common.RandomString(20)
	return user.currentToken
}

// CheckToken : Check user websocket authorization token
func (user *RealUser) CheckToken(token string) bool {
	return token == user.currentToken
}

// CurrentToken return user's current token
func (user *RealUser) CurrentToken() string {
	return user.currentToken
}

func (user *RealUser) ownsDevice(deviceID string) bool {
	return common.StringInSlice(deviceID, user.Devices)
}
