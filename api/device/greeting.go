package device

import (
	"encoding/json"

	"../../common"
	"github.com/gin-gonic/gin"
)

func greeting(c *gin.Context) {
	mdb, err := common.Mongo()
	defer mdb.Close()

	// check input
	deviceID := c.Query("id")

	// update
	var data gin.H
	col := mdb.DB("CUSmartFarm").C("devices")
	if common.PrintError(err) {
		c.JSON(500, "something went wrong")
		return
	}
	err = col.Find(gin.H{
		"id": deviceID,
	}).One(&data)
	if err != nil {
		c.JSON(404, "no device")
		return
	}
	payload, err := json.Marshal(gin.H{
		"t":      "cmd",
		"status": data["status"],
	})
	common.Println("greeting" + string(payload))
	common.PublishToMQTT([]byte("CUSmartFarm"), []byte(string(payload)))
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "greeting OK",
	})
}
