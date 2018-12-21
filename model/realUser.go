package model

import (
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

type RealUser struct {
	Username     string
	currentToken string
	devices      []string
	conn         *websocket.Conn
}

func (user *RealUser) Command(relay string, workmode string, payload interface{}) {

}
func (user *RealUser) ReportStatus(payload interface{}) {

}

func (user *RealUser) AddDevice(sensorID string, sensorInfo interface{}) {

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
