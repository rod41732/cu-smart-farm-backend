package user

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

// GetDevList : get user device list
func (user *RealUser) GetDevList() {
	resp, err := json.Marshal(gin.H{
		"t":       "response",
		"e":       "getDevList",
		"payload": user.devices,
	})
	common.PrintError(err)
	user.conn.WriteMessage(1, resp)
}
