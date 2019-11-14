package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/config"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/rod41732/cu-smart-farm-backend/router"
)

func main() {
	config.Init()
	common.InitializeKeyPair()
	middleware.Initialize()

	router.InitMQTT()
	go mqtt.MQTT()

	r := gin.Default()

	router.SetUpHTTPAPI(r)
	// seperated web socket service
	/*
		ws := r.Group("/subscribe")
		if common.Secure {
			ws.Use(middleware.UserAuth.MiddlewareFunc())
		}
		ws.GET("/ws", router.WebSocket)
	*/
	r.Run(":3000")

}
