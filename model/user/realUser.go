package user

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
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
func (user *RealUser) Init(devices []string, conn *websocket.Conn) {
	user.RegenerateToken()
	user.devices = devices
	user.conn = conn
}

// ReportStatus sends MQTT data to user via WebSocket
func (user *RealUser) ReportStatus(payload interface{}) {
	resp, _ := json.Marshal(payload)
	user.conn.WriteMessage(1, resp) // 1 is text message
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
