package main

import (
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"

	maniRouter "github.com/rod41732/cu-smart-farm-backend/router" 
	"github.com/rod41732/cu-smart-farm-backend/websocket-part/router" 

)

func main() {
	config.Init()
	common.InitializeKeyPair()
	middleware.Initialize()


	maniRouter.InitMQTTWithoutPublish() //s et handler eefore 
	go mqtt.MQTT()/// then connect and use that handler

	r := gin.Default()
	ws := r.Group("/subscribe")
	if common.Secure {
		ws.Use(middleware.UserAuth.MiddlewareFunc())
	}
	common.Println("[WS] Started successfully")
	ws.GET("/ws", router.WebSocket)

	r.Run(":3001")

}//////   
