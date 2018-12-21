package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/router"
)

func main() {

	common.InitializeKeyPair()

	r := gin.Default()
	middleware.Initialize()

	router.SetUpHttpAPI(r)
	r.Run(":3000")

}
