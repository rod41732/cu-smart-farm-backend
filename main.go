package main

import (
	"./api/middleware"
	"./common"
	"./router"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitializeKeyPair()

	go router.MQTT()
	r := gin.Default()
	middleware.Initialize()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world")
	})

	router.SetUpHttpAPI(r)
	r.Run(":3000")

}
