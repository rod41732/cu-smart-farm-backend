package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model/device"
	mMessage "github.com/rod41732/cu-smart-farm-backend/model/message"
	"github.com/rod41732/cu-smart-farm-backend/model/user"
	"gopkg.in/mgo.v2/bson"
)

var wsUpgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     wsCheckOrigin,
}

func wsCheckOrigin(r *http.Request) bool {
	return true // TODO origin check
}

// WebSocket : WS request handler
func WebSocket(c *gin.Context) {
	// get user part
	// headerUser, _ := c.Get("username")
	var headerUser interface{} = "rod41732" // TODO re-enable user check ing WS
	username := headerUser.(string)
	userObject := storage.GetUserStateInfo(username)
	wsRouter(c.Writer, c.Request, userObject)
}

func getDeviceAndParamFromMessage(payload map[string]interface{}) (device *device.Device, param map[string]interface{}, errmsg string) {
	var message mMessage.Message
	errmsg = ""
	if message.FromMap(payload) != nil {
		errmsg = "Bad Request"
	}
	param = message.Param
	device, err := storage.GetDevice(message.DeviceID)
	if err != nil {
		errmsg = "Device not found"
	}
	return
}

func responseStateBody(EndPoint string, success bool, errmsg string, nextToken string, data interface{}) []byte {
	result, err := json.Marshal(
		gin.H{
			"t":      "status",
			"e":      EndPoint,
			"status": gin.H{"success": success, "errmsg": errmsg},
			"token":  nextToken,
			"data":   data,
		})
	common.PrintError(err)
	return result
}

func wsRouter(w http.ResponseWriter, r *http.Request, clnt *user.RealUser) { // pass as pointer to prevent copying array
	common.Println("Web socket Connected")
	common.Println("Hi", clnt.Username)
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	clnt.SetConn(conn)

	resp, _ := json.Marshal(bson.M{"token": clnt.CurrentToken()}) // give client first token
	conn.WriteMessage(1, resp)

	if common.PrintError(err) {
		fmt.Println("  At WebSocket/wsRouter - Upgrading connection")
		return
	}
	for { // loop for command
		t, msg, err := conn.ReadMessage()
		if common.PrintError(err) {
			fmt.Println("  At WebSocket/wsRouter - Reading Message")
			break
		}

		var payload mMessage.APICall
		err = json.Unmarshal(msg, &payload)

		success := false
		hasGenToken := false
		errmsg := ""
		var data interface{}
		if err != nil {
			success, errmsg = false, "Bad Payload !!"
		} else {
			if !clnt.CheckToken(payload.Token) && false { // disable check : TODO re-enable check
				common.Println("[!] [WS] -- Invalid token")
				success, errmsg = false, "Invalid token"
			} else {
				common.Printf("[!] [WS] -- Endpoint : %s", payload.EndPoint)
				clnt.RegenerateToken()
				hasGenToken = true

				// Device Middleware
				var dev *device.Device
				var param map[string]interface{}
				dev, param, errmsg = getDeviceAndParamFromMessage(payload.Payload)
				if errmsg != "" {
					conn.WriteMessage(t, responseStateBody(payload.EndPoint, success, errmsg, clnt.CurrentToken(), nil))
					continue
				}
				switch payload.EndPoint {
				// case "addDevice":
				// 	success, errmsg = clnt.AddDevice(param, dev)
				// case "removeDevice":
				// 	success, errmsg = clnt.RemoveDevice(dev)
				case "pollDevice":
					dev.Poll()
					success, errmsg = true, "OK"
				case "setDevice":
					success, errmsg = clnt.SetDevice(param, dev)
				case "getLatestState":
					var results []client.Result
					success, errmsg, results = clnt.QueryDeviceLog(bson.M{"limit": 1}, dev)
					if len(results) > 0 && len(results[0].Series) > 0 && len(results[0].Series[0].Values) > 0 {
						data = zip(results[0].Series[0].Columns, results[0].Series[0].Values[0])
					}
					// case "getDevList":
				// 	clnt.GetDevList()
				// 	success, errmsg = true, "OK"
				default:
					success, errmsg = false, "unknown command"
				}
			}
		}
		var nextToken string
		if hasGenToken {
			nextToken = clnt.CurrentToken()
		} else {
			nextToken = ""
		}
		conn.WriteMessage(t, responseStateBody(payload.EndPoint, success, errmsg, nextToken, data))
		time.Sleep(time.Millisecond * 10)
	}
	clnt.SetConn(nil)
	common.Println("Web socket Disconnected")
}

func zip(keys []string, vals []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i, k := range keys {
		result[k] = vals[i]
	}
	return result
}
