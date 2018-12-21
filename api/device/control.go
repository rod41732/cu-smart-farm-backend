package device

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

var num2state = map[string]string{
	"0": "OFF",
	"1": "ON",
	"2": "AUTO",
}

func setState(c *gin.Context) {
	mdb, err := common.Mongo()
	defer mdb.Close()

	// check input
	relay := c.Query("relay")
	x, err := strconv.Atoi(relay)
	if err != nil || !(1 <= x && x <= 4) { // 1-4 = each device, 5 = all device
		c.JSON(400, "bad relay")
		return
	}
	state, ok := num2state[c.Query("state")]
	if !ok {
		c.JSON(400, "bad state")
	}
	deviceID := c.Query("id")

	// update
	var data, status map[string]interface{}
	col := mdb.DB("CUSmartFarm").C("devices")
	if common.PrintError(err) {
		c.JSON(500, "something went wrong")
		return
	}
	err = col.Find(gin.H{
		"id": deviceID,
	}).One(&data)
	// modify status
	status = data["status"].(map[string]interface{})
	status[relay] = state

	err = col.Update(gin.H{
		"id": deviceID,
	}, gin.H{
		"$set": gin.H{"status": status},
	})
	if common.PrintError(err) {
		c.JSON(500, err)
		return
	}

	status = map[string]interface{}{"t": "cmd", "cmd": "set", "state": status}
	payload, err := json.Marshal(status)
	common.PublishToMQTT([]byte("CUSmartFarm"), []byte(payload))
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "sent " + string(payload) + " to MQTT update",
	})
}

func fetchInfo(c *gin.Context) {
	payload, err := json.Marshal(map[string]interface{}{
		"t":   "cmd",
		"cmd": "fetch",
	})
	if common.PrintError(err) {
		c.JSON(500, "server error")
		return
	}

	common.PublishToMQTT([]byte("CUSmartFarm"), []byte(payload))
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "sent " + string(payload) + " to MQTT",
	})
}
