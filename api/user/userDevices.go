package user

import (
	"encoding/json"

	"github.com/rod41732/cu-smart-farm-backend/model/message"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/storage"
)

// shortcut to send 500 error
func error500(c *gin.Context) {
	c.JSON(500, gin.H{
		"success": false,
		"message": "Something went wrong",
	})
}

func addDevice(c *gin.Context) {
	usr, _ := c.Get("user")
	user, _ := usr.(*middleware.User)
	userObject := storage.GetUserStateInfo(user.Username)
	if userObject == nil {
		error500(c)
		return
	}

	var payload message.Message
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err != nil {
		common.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Payload format",
		})
	} else {
		dev, err := storage.GetDevice(payload.DeviceID)
		common.PrintError(err)
		var ok bool
		var errmsg string
		if err != nil {
			ok, errmsg = false, "Invalid device"
		} else {
			ok, errmsg = userObject.AddDevice(payload.Param, dev)
		}
		var status int
		if !ok { // TODO: spaghetti
			status = 500
		} else {
			status = 200
		}
		c.JSON(status, gin.H{
			"success": ok,
			"message": errmsg,
		})
	}
}

func removeDevice(c *gin.Context) {
	usr, _ := c.Get("user")
	user, _ := usr.(*middleware.User)
	userObject := storage.GetUserStateInfo(user.Username)
	if userObject == nil {
		error500(c)
		return
	}

	var payload message.Message
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Request",
		})
	} else {
		dev, err := storage.GetDevice(payload.DeviceID)
		var ok bool
		var errmsg string
		// var errmsg string
		if err != nil {
			ok, errmsg = false, "GetDevice not found"
		} else {
			ok, errmsg = userObject.RemoveDevice(dev)
		}
		var status int
		if !ok { // TODO: spaghetti
			status = 500
		} else {
			status = 200
		}
		c.JSON(status, gin.H{
			"success": ok,
			"message": errmsg,
		})
	}
}

func setDevice(c *gin.Context) {
	usr, _ := c.Get("user")
	user, _ := usr.(*middleware.User)
	userObject := storage.GetUserStateInfo(user.Username)
	if userObject == nil {
		error500(c)
		return
	}

	var payload message.Message
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err == nil {
		common.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Request",
		})
	} else {
		dev, err := storage.GetDevice(payload.DeviceID)
		var ok bool
		var errmsg string
		if err != nil {
			ok, errmsg = false, "Device not found"
		} else {
			ok, errmsg = userObject.SetDevice(payload.Param, dev)
		}
		var status int
		if !ok { // TODO: spaghetti
			status = 500
		} else {
			status = 200
		}
		c.JSON(status, gin.H{
			"success": ok,
			"message": errmsg,
		})
	}
}

func getDevicesList(c *gin.Context) {
	usr, _ := c.Get("user")
	user, _ := usr.(*middleware.User)
	userObject := storage.GetUserStateInfo(user.Username)
	if userObject == nil {
		error500(c)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "OK",
		"data":    userObject.Device  s(),
	})
}
