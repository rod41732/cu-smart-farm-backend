package user

import (
	"encoding/json"
	"net/http"

	"github.com/rod41732/cu-smart-farm-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     wsCheckOrigin,
}

type userDevice struct {
	ownedDevices []string
}

// WebSocket : WebSocket request handling
func WebSocket(c *gin.Context) {
	var userDev userDevice
	userSession, _ := c.Get("username")
	username := userSession.(*middleware.User).Username
	mdb, err := common.Mongo()

	if common.PrintError(err) {
		c.JSON(500, gin.H{
			"msg": "Connection to database failed",
		})
		return
	}

	collection := mdb.DB("CUSmartFarm").C("devices")
	collection.Find(bson.M{
		"username": username,
	}).One(&userDev)

	clientState := model.RealUser{
		Username: username,
	}
	// User Authorization
	wsRouter(c.Writer, c.Request)
	// common.NewWebsocketConnectionInfo(userSession, conn)
}

func wsRouter(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if common.PrintError(err) {
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if common.PrintError(err) {
			break
		}
		println("Message from client webSocket:", string(msg))

		incomeCommand := common.WsCommand{}
		err = json.Unmarshal(msg, &incomeCommand)
		if !common.PrintError(err) {
			switch incomeCommand.Endpoint {
			case "addDevice":
			default:
			}

		}
		wsReponseStatus(incomeCommand.Endpoint, err, conn, t)
	}
}

func wsCheckOrigin(r *http.Request) bool {
	return true
}

func wsReponseStatus(endpoint string, err error, conn *websocket.Conn, msgType int) {
	callbackMessage := common.WsCommand{
		Endpoint: endpoint,
		Payload:  "success",
	}
	if err != nil {
		callbackMessage.Payload = err.Error()
	}
	jsonMsg, _ := json.Marshal(callbackMessage)
	conn.WriteMessage(msgType, jsonMsg)
}
