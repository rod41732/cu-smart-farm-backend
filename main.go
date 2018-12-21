package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/rod41732/cu-smart-farm-backend/router"
)

func main() {

	common.InitializeKeyPair()
	middleware.Initialize()

	go mqtt.MQTT()

	r := gin.Default()

	router.SetUpHttpAPI(r)
	r.GET("/ws", router.WebSocket)
	r.Run(":3000")

}
