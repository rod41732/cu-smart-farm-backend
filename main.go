package main

import (
	"./router"
	"github.com/gin-gonic/gin"
)

func main() {

	go router.MQTT()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world")
	})

	router.SetUpHttpAPI(r)
	r.Run(":3000")

}
