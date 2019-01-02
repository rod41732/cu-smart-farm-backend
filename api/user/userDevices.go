package user

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	User "github.com/rod41732/cu-smart-farm-backend/model/user"
	"github.com/rod41732/cu-smart-farm-backend/storage"
)

// shortcut to send 500 error
func error500(c *gin.Context) {
	c.JSON(500, gin.H{
		"success": false,
		"message": "some errorr ",
	})
}

func addDevice(c *gin.Context) {
	usr, ok := c.Get("user")
	if !ok {
		error500(c)
		return
	}
	user, ok := usr.(*middleware.User)
	if !ok {
		error500(c)
		return
	}
	userObject, ok := storage.GetUserStateInfo(user.Username).(*User.RealUser)
	if !ok {
		error500(c)
		return
	}

	var payload gin.H
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err != nil {
		common.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Request",
		})
	} else {
		ok, errmsg := userObject.AddDevice(payload)
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
	usr, ok := c.Get("user")
	if !ok {
		error500(c)
		return
	}
	user, ok := usr.(*middleware.User)
	if !ok {
		error500(c)
		return
	}
	userObject, ok := storage.GetUserStateInfo(user.Username).(*User.RealUser)
	if !ok {
		error500(c)
		return
	}
	var payload gin.H
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Request",
		})
	} else {
		ok, errmsg := userObject.RemoveDevice(payload)
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
	usr, ok := c.Get("user")
	if !ok {
		error500(c)
		return
	}
	user, ok := usr.(*middleware.User)
	if !ok {
		error500(c)
		return
	}
	userObject, ok := storage.GetUserStateInfo(user.Username).(*User.RealUser)
	if !ok {
		error500(c)
		return
	}
	var payload gin.H
	err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)

	if err != nil {
		common.Println(err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "Bad Request",
		})
	} else {
		ok, errmsg := userObject.SetDevice(payload)
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
	usr, ok := c.Get("user")
	if !ok {
		error500(c)
		return
	}
	user, ok := usr.(*middleware.User)
	if !ok {
		error500(c)
		return
	}
	userObject, ok := storage.GetUserStateInfo(user.Username).(*User.RealUser)
	if !ok {
		error500(c)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "OK",
		"data":    userObject.Devices(),
	})
}
