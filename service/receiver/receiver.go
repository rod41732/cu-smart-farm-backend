package receiver

import (
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/service/storage"
)

type Trigger string

// Update triggers storage:LoadDevice
func (t *Trigger) Update(deviceID string, errmsg *string) error {
	common.Println("[Trigger] received trigger", deviceID)
	*errmsg = "OK"
	err := storage.LoadDevice(deviceID)
	if err != nil {
		*errmsg = "LoadDevice Error"
		return err
	}
	return nil
}
