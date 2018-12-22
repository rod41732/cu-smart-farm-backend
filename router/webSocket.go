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
	// tmp, _ := c.Get("username")
	// user := tmp.(middleware.User)
	user := middleware.User{Username: "rod41732"}
	wsRouter(c.Writer, c.Request, user.Username)
}

type userData struct {
	Devices []string `json:"devices"`
}

func wsRouter(w http.ResponseWriter, r *http.Request, username string) {
	common.Println("Web socket Connected")
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return
	}
	var temp userData
	mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": username,
	}).One(&temp)

	client := user.RealUser{Username: username}
	client.Init(temp.Devices, conn)
	storage.SetUserStateInfo(username, &client)

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

		var success bool
		var errmsg string
		if err != nil {
			common.PrintError(err)
			success, errmsg = false, "No Command Specified"
		} else {
			switch payload.EndPoint {
			case "addDevice":
				success, errmsg = client.AddDevice(payload.Payload)
			case "removeDevice":
				success, errmsg = client.RemoveDevice(payload.Payload)
			case "pollDevice":
				success, errmsg = client.PollDevice(payload.Payload)
			default:
				success, errmsg = false, "unknown command"
			}
		}

		resp, err := json.Marshal(bson.M{"success": success, "errmsg": errmsg})
		conn.WriteMessage(t, resp)
		time.Sleep(time.Millisecond * 10)
	}
	storage.SetUserStateInfo(username, &user.NullUser{})
	common.Println("Web socket Disconnected")
}
