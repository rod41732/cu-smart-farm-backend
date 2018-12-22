package user

import "github.com/rod41732/cu-smart-farm-backend/common"

// NullUser is "placeholder" when client is disconnected
type NullUser struct {
}

// ReportStatus for null user do nothing
func (user *NullUser) ReportStatus(payload interface{}) {
	common.Printf("received: %v\n", payload)
	// return
}
