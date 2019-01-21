package user

import (
	"encoding/json"
	"errors"

	"github.com/rod41732/cu-smart-farm-backend/common"

	"github.com/influxdata/influxdb/client/v2"

	"github.com/rod41732/cu-smart-farm-backend/model/device"
	"gopkg.in/mgo.v2/bson"

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

// return userObject based on name stored in gin.Context
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

// return current username
func checkStatus(c *gin.Context) {
	username, _ := c.Get("user")
	c.JSON(200, bson.M{
		"username": username,
	})
}

// extract param from payload (post body), get device from storage
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

func setRelay(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.SetDeviceRelay(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func setRelayName(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.SetDeviceRelayName(param, dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
	})
}

func getDevicesList(c *gin.Context) {
	type deviceShortInfo struct {
		Name        string `json:"name"`
		ID          string `json:"deviceID"`
		Description string `json:"desc"`
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
			devShortInfo[i].Description = devInfo.Description
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    devShortInfo,
	})
}

// TODO: change error repsonse
func editDevice(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg = user.EditDevice(param, dev)
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

func getRelayNames(c *gin.Context) {
	ok, errmsg := true, "OK"
	var names []string
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, _, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg, names = user.GetRelayNames(dev)
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    names,
	})
}

func getDeviceLog(c *gin.Context) {
	ok, errmsg := true, "OK"
	var log interface{}
	var results []client.Result
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		dev, param, err := extractDeviceIDandParam(c)
		if err != nil {
			ok, errmsg = false, err.Error()
		} else {
			ok, errmsg, results = user.QueryDeviceLog(param, dev)
			if common.HaveSeries(results) {
				log = results[0].Series[0]
			}
		}
	}

	c.JSON(status(ok), gin.H{
		"success": ok,
		"message": errmsg,
		"data":    log,
	})
}
