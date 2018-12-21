package user

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     wsCheckOrigin,
}

// WebSocket : WebSocket request handling
func WebSocket(c *gin.Context) {
	wsRouter(c.Writer, c.Request, sessions.Default(c))
}

func wsRouter(w http.ResponseWriter, r *http.Request, userSession sessions.Session) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if common.CheckErr(err) {
		return
	}
	common.NewWebsocketConnectionInfo(userSession, conn)
	previousMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
	go func() {
		for {
			// ****************** Device Log websocket broadcast simulator ******************
			if (time.Now().UnixNano()/int64(time.Millisecond))-previousMillisecond > 2000 {
				mockDeviceLog := gin.H{
					"aX":     rand.Float64() * 10,
					"aY":     rand.Float64() * 10,
					"aZ":     rand.Float64() * 10,
					"angleX": rand.Float64() * 90,
					"angleY": rand.Float64() * 90,
					"angleZ": rand.Float64() * 90,
					"gyroX":  rand.Float64() * 10,
					"gyroY":  rand.Float64() * 10,
					"gyroZ":  rand.Float64() * 10,
				}
				jsonMsg, _ := json.Marshal(mockDeviceLog)
				common.BroadCastAll("deviceLog", string(jsonMsg))
				println("simulated deviceLog")
				previousMillisecond = (time.Now().UnixNano() / int64(time.Millisecond))
			}
		}
	}()
	for {
		t, msg, err := conn.ReadMessage()
		if common.CheckErr(err) {
			break
		}
		println("Message from client webSocket:", string(msg))

		incomeCommand := common.WsCommand{}
		err = json.Unmarshal(msg, &incomeCommand)
		if !common.CheckErr(err) {
			switch incomeCommand.Event {
			case "throttle":
				err = Device.SetThrottle(incomeCommand.Payload)
			case "DeviceMotorSpeed":
				err = Device.SetMotor(incomeCommand.Payload)
			case "TurnDeviceCamLeft":
				err = Device.SetCamera("Left", incomeCommand.Payload)
			case "TurnDeviceCamRight":
				err = Device.SetCamera("Right", incomeCommand.Payload)
			case "GetDeviceStatus":
				common.BroadCastAll("DeviceStatus", Client.EarliestDevicePhysicalLog("HardCode101"))
			}

		}
		wsReponseStatus(incomeCommand.Event, err, conn, t)
	}
}

func wsCheckOrigin(r *http.Request) bool {
	return true
}

func wsReponseStatus(event string, err error, conn *websocket.Conn, msgType int) {
	callbackMessage := common.WsCommand{
		Event:   event,
		Payload: "success",
	}
	if err != nil {
		callbackMessage.Payload = err.Error()
	}
	jsonMsg, _ := json.Marshal(callbackMessage)
	conn.WriteMessage(msgType, jsonMsg)
}
