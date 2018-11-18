package middleware

import (
	"../../common"
	"github.com/gin-gonic/gin"
)

func OwnerCheck(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.Resp(500, err, c) {
		c.Abort()
		return
	}
	deviceId := c.Query("deviceId")
	var device *gin.H
	col := mdb.DB("CUSmartFarm").C("device")
	col.Find(gin.H{
		"id": deviceId,
	}).One(device)

	v, found := c.Get("username")
	if !found {
		c.JSON(403, "Unauthorized")
		c.Abort()
		return
	}
	user, ok := v.(*User)
	if !ok || user == nil {
		c.JSON(403, "Unauthorized")
		c.Abort()
		return
	}
	common.Printf("============ %#v , %#v\n", user, deviceId)
	common.Println(user.Username == deviceId)
	if user.Username == deviceId {
		c.Next()
	} else {
		c.JSON(403, "Not your device")
		c.Abort()
		return
	}
}
