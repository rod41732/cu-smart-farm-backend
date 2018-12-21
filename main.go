package main

import (
	"./api/middleware"
	"./common"
	"./router"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitializeKeyPair()
	middleware.Initialize()

	// go router.MQTT()
	r := gin.Default()

	router.SetUpHttpAPI(r)
	r.GET("/ws", router.WebSocket)
	r.Run(":3000")

}
