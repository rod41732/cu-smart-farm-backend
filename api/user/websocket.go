package user

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     wsCheckOrigin,
}

// WebSocket : WebSocket request handling
func WebSocket(c *gin.Context) {
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
