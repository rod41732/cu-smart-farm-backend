package user

/*
import (
	"encoding/json"
	"net/http"

	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2/bson"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     wsCheckOrigin,
}

type userData struct {
	username     string
	ownedDevices []string
}

// WebSocket : WebSocket request handling
func WebSocket(c *gin.Context) {
	var userInfo userData
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
	}).One(&userInfo)
	for _, deviceID := range userInfo.ownedDevices {
		mqtt.SubscribeDevice(deviceID)
	}
	// User Authorization
	wsRouter(c.Writer, c.Request, userInfo)
}

func wsRouter(w http.ResponseWriter, r *http.Request, userInfo userData) {
	client := model.RealUser{
		Username: userInfo.username,
	}
	conn, err := wsupgrader.Upgrade(w, r, nil)
	client.Init(userInfo.ownedDevices, conn)
	storage.SetUserStateInfo(userInfo.username, model.User(client))
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
*/
