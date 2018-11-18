package device

import (
	"encoding/json"
	"strconv"

	"../../common"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

var num2state = map[string]string{
	"0": "OFF",
	"1": "ON",
	"2": "AU",
}

func setState(c *gin.Context) {
	mdb, err := common.Mongo()
	defer mdb.Close()

	// check input
	relay := c.Query("relay")
	x, err := strconv.Atoi(relay)
	if err != nil || !(1 <= x && x <= 5) { // 1-4 = each device, 5 = all device
		c.JSON(400, "bad relay")
		return
	}
	state, ok := num2state[c.Query("state")]
	common.Println("============== state = ", state)
	if !ok {
		c.JSON(400, "bad state")
	}
	deviceID := c.Query("deviceId")

	// update
	var data, status map[string]interface{}
	col := mdb.DB("CUSmartFarm").C("devices")
	if common.PrintError(err) {
		c.JSON(500, "something went wrong")
		return
	}
	err = col.Find(bson.M{
		"id": deviceID,
	}).One(&data)
	if err != nil {
		c.JSON(404, "no device")
		return
	}

	// modify status
	status = data["status"].(map[string]interface{})
	status[relay] = state

	err = col.Update(bson.M{
		"id": deviceID,
	}, bson.M{
		"$set": bson.M{"status": status},
	})
	if common.PrintError(err) {
		c.JSON(500, err)
		return
	}
	status = map[string]interface{}{"t": "cmd", "status": status}
	payload, err := json.Marshal(status)
	common.PublishToMQTT([]byte("CUSmartFarm"), []byte(payload))
	c.JSON(200, gin.H{
		"success": true,
		"msg":     "sent " + string(payload) + " to MQTT update",
	})
}
