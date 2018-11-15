package device

import (
	"encoding/json"
	"fmt"
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

// DeviceControlAPI : sets up device control API
func DeviceControlAPI(r *gin.RouterGroup) {
	deviceAPI := r.Group("/device")

	deviceAPI.GET("/set", func(c *gin.Context) {
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
		deviceId := c.Query("id")

		// update
		var data, status map[string]interface{}
		col := mdb.DB("CUSmartFarm").C("devices")
		if common.CheckErr("connect to mongo", err) {
			c.JSON(500, "something went wrong")
			return
		}
		err = col.Find(bson.M{
			"id": deviceId,
		}).One(&data)
		if err != nil {
			c.JSON(404, "no device")
			return
		}

		// modify status
		status = data["status"].(map[string]interface{})
		status["relay"+c.Query("1234")] = state

		err = col.Update(bson.M{
			"id": deviceId,
		}, bson.M{
			"$set": bson.M{"status": status},
		})
		if common.CheckErr("update", err) {
			c.JSON(500, err)
			return
		}

		payload := state + relay
		common.PublishToMQTT([]byte("CUSmartFarm"), []byte(payload))
		common.Println(data)
		c.JSON(200, gin.H{
			"success": true,
			"msg":     "sent " + payload + " to MQTT update",
		})
	})

	deviceAPI.GET("/greeting", func(c *gin.Context) {
		mdb, err := common.Mongo()
		defer mdb.Close()

		// check input
		deviceId := c.Query("id")

		// update
		var data map[string]interface{}
		col := mdb.DB("CUSmartFarm").C("devices")
		if common.CheckErr("connect to mongo", err) {
			c.JSON(500, "something went wrong")
			return
		}
		err = col.Find(bson.M{
			"id": deviceId,
		}).One(&data)
		if err != nil {
			c.JSON(404, "no device")
			return
		}
		common.Println("-------", data["status"])
		payload, err := json.Marshal(data["status"])
		spayload := fmt.Sprintf("%s", payload)
		common.PublishToMQTT([]byte("CUSmartFarm"), payload)
		common.Println("======= sending", spayload)
		c.JSON(200, gin.H{
			"success": true,
			"msg":     spayload,
		})
	})
}
