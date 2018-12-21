package main

import (
	"./api/middleware"
	"./common"
	"./router"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitializeKeyPair()

	r := gin.Default()
	middleware.Initialize()

	router.SetUpHttpAPI(r)
	r.Run(":3000")

}
