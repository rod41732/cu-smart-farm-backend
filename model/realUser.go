package model

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type RealUser struct {
	Username     string
	currentToken string
	devices      []string
	conn         *websocket.Conn
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (user *RealUser) Init(devices []string, conn *websocket.Conn) {
	user.RegenerateToken()
	user.devices = devices
	user.conn = conn
}

func (user *RealUser) Command(relay string, workmode string, payload interface{}) {

}
func (user *RealUser) ReportStatus(payload interface{}) {

}

func (user *RealUser) AddDevice(sensorID string, sensorInfo interface{}) {

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
