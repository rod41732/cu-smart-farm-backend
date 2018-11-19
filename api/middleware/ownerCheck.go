package middleware

import (
	"../../common"
	"github.com/gin-gonic/gin"
)

// OwnerCheck this middleware check whether `id` in query match `Username`
func OwnerCheck(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.Resp(500, err, c) {
		c.Abort()
		return
	}
	deviceId := c.Query("id")
	var device map[string]interface{}
	col := mdb.DB("CUSmartFarm").C("devices")
	col.Find(gin.H{
		"id": deviceId,
	}).One(&device)

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

	common.Printf("============ user %#v , %#v %#v\n", user.Username, device["owner"], device)
	// common.Println(user.Username == device["owner"])
	if user.Username == device["owner"] {
		// if true {
		common.Println("ownerchecko ok")
		c.Next()
	} else {
		c.JSON(403, "Not your device")
		c.Abort()
		return
	}
}
