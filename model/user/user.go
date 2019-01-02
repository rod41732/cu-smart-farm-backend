package user

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
)

// RealUser represent client connected via WebSocket
type RealUser struct {
	Username     string
	currentToken string
	devices      []string
	conn         *websocket.Conn
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//Init  initializes user
func (user *RealUser) Init(devices []string) {
	user.RegenerateToken()
	user.devices = devices
	user.conn = nil
}

//SetConn is setter for Conn
func (user *RealUser) SetConn(conn *websocket.Conn) {
	user.conn = conn
}

// ReportStatus sends MQTT data to user via WebSocket then insert into InfluxDB
func (user *RealUser) ReportStatus(payload model.DeviceMessagePayload, deviceID string) {
	resp, _ := json.Marshal(payload)
	if user.conn != nil {
		user.conn.WriteMessage(1, resp) // 1 is text message
	} else {
		common.Println("[User] null user => ", resp)
	}
	common.WriteInfluxDB("cu_smartfarm_sensor_log", map[string]string{"device": deviceID}, payload.ToMap())
}

// RegenerateToken : Regenerate user websocket authorization token
func (user *RealUser) RegenerateToken() string {
	user.currentToken = randomString(20)
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

// randomString : helper function for random string with custom length and charset
func randomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (user *RealUser) ownsDevice(deviceID string) bool {
	return common.StringInSlice(deviceID, user.devices)
}

func (user *RealUser) Devices() []string {
	return user.devices
}
