package user

import (
	"encoding/json"
	"errors"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/model/message"
	"github.com/rod41732/cu-smart-farm-backend/model/user"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/storage"
)

// return 200 if OK else 500
func status(ok bool) int {
	if ok {
		return 200
	}
	return 500
}

func extractUser(c *gin.Context) (userObject *user.RealUser, err error) {
	usr, ok := c.Get("user")
	if !ok {
		err = errors.New("Not logged in")
		return
	}
	username, ok := usr.(string)
	if !ok {
		err = errors.New("Something went wrong, please login again")
		return
	}
	userObject = storage.GetUserStateInfo(username)
	return
}

func extractDeviceIDandParam(c *gin.Context) (dev *device.Device, param map[string]interface{}, err error) {
	var payload message.Message
	err = json.Unmarshal([]byte(c.PostForm("payload")), &payload)
	if err != nil {
		err = errors.New("Bad payload")
		return
	}
	param = payload.Param
	dev, err = storage.GetDevice(payload.DeviceID)
	return
}

// shortcut to send 500 error
func error500(c *gin.Context) {
	c.JSON(500, gin.H{
		"success": false,
		"message": "Something went wrong",
	})
}

func addDevice(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.AddDevice(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func removeDevice(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, _, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.RemoveDevice(dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func setDevice(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.SetDevice(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func getDevicesList(c *gin.Context) {
	type deviceShortInfo struct {
		Name string `json:"name"`
		ID   string `json:"deviceID"`
	}
	var devices []string
	var devShortInfo []deviceShortInfo

	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		devices = user.Devices
		devShortInfo = make([]deviceShortInfo, len(devices))
		for i, device := range devices {
			devInfo, _ := storage.GetDevice(device)
			devShortInfo[i].ID = devInfo.ID
			devShortInfo[i].Name = devInfo.Name
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    devShortInfo,
	})
}

// TODO: change error repsonse
func renameDevice(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.RenameDevice(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func getDeviceInfo(c *gin.Context) {
	ok, errmsg := true, "OK"
	var info interface{}
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, _, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg, info = user.GetDeviceInfo(dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    info,
	})
}

func getDeviceLog(c *gin.Context) {
	ok, errmsg := true, "OK"
	var log interface{}
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg, log = user.QueryDeviceLog(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    log,
	})
}
