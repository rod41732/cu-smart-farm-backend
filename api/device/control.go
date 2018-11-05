package device

import (
	"encoding/json"
	"strconv"

	"../../common"
	"github.com/gin-gonic/gin"
	"github.com/surgemq/message"
)

// DeviceControlAPI : sets up device control API
func DeviceControlAPI(r *gin.RouterGroup) {
	deviceAPI := r.Group("/device")

	deviceAPI.GET("/set", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
		if err != nil || !(1 <= id && id <= 5) { // 1-4 = each device, 5 = all device
			c.JSON(400, gin.H{
				"success": false,
				"msg":     "id must be 1 to 5",
			})
		}
		state, err := strconv.Atoi(c.DefaultQuery("state", "-1"))
		if err != nil || !(0 <= state && state <= 3) {
			c.JSON(400, gin.H{
				"success": false,
				"msg":     "state must be 0 to 3",
			})
		}
		// dpn't set to 2 because SmartFarm Board doesn't support QoS level 2
		x := map[int]string{
			0: "OFF",
			1: "ON",
			2: "MA", // Mode Auto (Value Threshold)
			3: "MM", // Mode Manual (Time)
		}

		// payload is in for <State><relay_number>
		payload := x[state]
		if id < 5 {
			payload += strconv.Itoa(id)
		} else {
			payload += "ALL"
		}

		// payload2 := map[string]string{
		// 	"command": payload,
		// 	"type":    "test",
		// }

		msg := message.NewPublishMessage()
		msg.SetTopic([]byte("CUSmartFarm"))
		msg.SetQoS(0)
		text, _ := json.Marshal(payload)
		msg.SetPayload([]byte(text))
		common.MqttClient.Publish(msg, nil)

		c.JSON(200, gin.H{
			"success": true,
			"msg":     "sent " + payload + " to MQTT",
		})

	})
}
