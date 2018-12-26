package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
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

	// db part
	mdb, err := common.Mongo()
	var dbUser userData
	defer mdb.Close()
	if common.PrintError(err) {
		c.JSON(500, gin.H{
			"msg": "Connection to database failed",
		})
		return
	}
	common.Printf("finding %s\n", username)
	mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": username,
	}).One(&dbUser)

	common.Printf("user = %#v\n", dbUser)
	// for _, deviceID := range userInfo.ownedDevices {
	// 	mqtt.SubscribeDevice(deviceID)
	// }
	wsRouter(c.Writer, c.Request, &dbUser)
}

type userData struct {
	Username string   `json:"username"`
	Devices  []string `json:"devices"`
}

func wsRouter(w http.ResponseWriter, r *http.Request, dbUser *userData) { // pass as pointer to prevent copying array
	common.Println("Web socket Connected")
	common.Println("Hi,", dbUser.Username)
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	username := dbUser.Username
	client := user.RealUser{Username: username}
	client.Init(dbUser.Devices, conn)
	storage.SetUserStateInfo(username, &client)

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
		if err != nil {
			common.PrintError(err)
			success, errmsg = false, "Bad Payload"
		} else {
			if !client.CheckToken(payload.Token) && false { // disable check
				success, errmsg = false, "Invalid token"
			} else {
				client.RegenerateToken()
				hasGenToken = true
				switch payload.EndPoint {
				case "addDevice":
					success, errmsg = client.AddDevice(payload.Payload)
				case "removeDevice":
					success, errmsg = client.RemoveDevice(payload.Payload)
				case "pollDevice":
					success, errmsg = client.PollDevice(payload.Payload)
				case "setDevice":
					success, errmsg = client.SetDevice(payload.Payload)
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
		resp, err := json.Marshal(bson.M{"success": success, "errmsg": errmsg, "token": nextToken})
		conn.WriteMessage(t, resp)
		time.Sleep(time.Millisecond * 10)
	}
	storage.SetUserStateInfo(username, &user.NullUser{})
	common.Println("Web socket Disconnected")
}
