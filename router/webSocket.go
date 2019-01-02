package router

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
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
	return true // todo origin check
}

// WebSocket : WS request handler
func WebSocket(c *gin.Context) {
	// get user part
	// tmp, _ := c.Get("username")
	// headerUser := tmp.(middleware.User)
	headerUser := middleware.User{Username: "rod41732"}
	username := headerUser.Username
	userObject := storage.GetUserStateInfo(username).(*user.RealUser)
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
	if common.PrintError(err) {
		errmsg = "Device not found"
	}
	return
}

func responseStateBody(EndPoint string, success bool, errmsg string, nextToken string) []byte {
	result, err := json.Marshal(
		gin.H{
			"t":      "status",
			"e":      EndPoint,
			"status": gin.H{"success": success, "errmsg": errmsg},
			"token":  nextToken,
		})
	common.PrintError(err)
	return result
}

func wsRouter(w http.ResponseWriter, r *http.Request, client *user.RealUser) { // pass as pointer to prevent copying array
	common.Println("Web socket Connected")
	common.Println("Hi", client.Username)
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	client.SetConn(conn)

	resp, _ := json.Marshal(bson.M{"token": client.CurrentToken()}) // give client first token
	conn.WriteMessage(1, resp)

	if common.PrintError(err) {
		return
	}
	for { // loop for command
		t, msg, err := conn.ReadMessage()
		if err != nil {
			common.PrintError(err)
			break
		}

		var payload mMessage.APICall
		err = json.Unmarshal(msg, &payload)

		success := false
		hasGenToken := false
		errmsg := ""
		if common.PrintError(err) {
			common.Println("[!] [WS] -- Bad Payload")
			success, errmsg = false, "Bad Payload"
		} else {
			if !client.CheckToken(payload.Token) && false { // disable check
				common.Println("[!] [WS] -- Invalid token")
				success, errmsg = false, "Invalid token"
			} else {
				common.Printf("[!] [WS] -- Endpoint : %s", payload.EndPoint)
				client.RegenerateToken()
				hasGenToken = true

				// Device Middleware
				var dev *device.Device
				var param map[string]interface{}
				if regexp.MustCompile("Device").MatchString(payload.EndPoint) {
					dev, param, errmsg = getDeviceAndParamFromMessage(payload.Payload)
					if errmsg != "" {
						conn.WriteMessage(t, responseStateBody(payload.EndPoint, success, errmsg, client.CurrentToken()))
						continue
					}
				}

				switch payload.EndPoint {
				case "addDevice":
					success, errmsg = client.AddDevice(param, dev)
				case "removeDevice":
					success, errmsg = client.RemoveDevice(dev)
				case "pollDevice":
					dev.Poll()
					success, errmsg = true, "OK"
				case "setDevice":
					success, errmsg = client.SetDevice(param, dev)
				case "getDevList":
					client.GetDevList()
					success, errmsg = true, "OK"
				default:
					success, errmsg = false, "unknown command"
				}
			}
		}
		var nextToken string
		if hasGenToken {
			nextToken = client.CurrentToken()
		} else {
			nextToken = ""
		}
		conn.WriteMessage(t, responseStateBody(payload.EndPoint, success, errmsg, nextToken))
		time.Sleep(time.Millisecond * 10)
	}
	client.SetConn(nil)
	common.Println("Web socket Disconnected")
}
