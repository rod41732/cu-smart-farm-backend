package router

import (
	"encoding/json"
	"net/http"
	"time"

	"../api/middleware"
	"../common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	sess := sessions.Default(c)
	user, _ := c.Get("username").(middleware.User)
	deviceId, ok := c.Get("id")
	sess.set("username", user.Username)
	sess.set("deviceId", deviceId)
	wsRouter(c.Writer, c.Request, sess)
}

func wsRouter(w http.ResponseWriter, r *http.Request, userSession sessions.Session) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if common.PrintError(err) {
		return
	}
	for { // loop for command
		t, msg, err := conn.ReadMessage()
		if common.PrintError(err) {
			break
		}
		var payload map[string]interface{}
		err := json.Unmarshal(msg, &payload)
		if !err {
			command, err := payload["cmd"].(string)
			if !err && command == "fetch" {
				if common.TellDevice(userSession.get("deviceId").(string)) {
					conn.WriteMessage(t, "OK");
				}
				else {
					conn.WriteMessage(t, "Error");
				}
			}
		}
		time.Sleep(time.Millisecond * 10)

	}

}
func deviceWebSocket(c *gin.Context) {
	sess := sessions.Default(c)
	deviceId, ok := c.Get("id")
	sess.set("deviceId", deviceId)
	deviceWsRouter(c.Writer, c.Request, sess)
}

func deviceWsRouter(w http.ResponseWriter, r *http.Request, userSession sessions.Session) {
	userSession.Set()
}
