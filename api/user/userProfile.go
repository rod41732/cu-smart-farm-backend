package user

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"
)

func editProfile(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)
		if err != nil {
			common.Printf("Payload is %#v\n", payload)
			ok, errmsg = false, "Bad Payload"
		} else {
			ok, errmsg = user.EditProfile(payload)
		}
	}

	c.JSON(status(ok), bson.M{
		"success": ok,
		"message": errmsg,
	})
}

func changePassword(c *gin.Context) {
	ok, errmsg := true, "OK"
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)
		if err != nil {
			ok, errmsg = false, "Bad Payload"
		} else {
			ok, errmsg = user.ChangePassword(payload)
		}
	}

	c.JSON(status(ok), bson.M{
		"success": ok,
		"message": errmsg,
	})
}

func getProfile(c *gin.Context) {
	ok, errmsg := true, "OK"
	var result map[string]interface{}
	user, err := extractUser(c)
	if err != nil {
		ok, errmsg = false, err.Error()
	} else {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(c.PostForm("payload")), &payload)
		if err != nil {
			ok, errmsg = false, "Bad Payload"
		} else {
			ok, errmsg, result = user.GetProfile()
		}
	}

	c.JSON(status(ok), bson.M{
		"success": ok,
		"message": errmsg,
		"data":    result,
	})
}
