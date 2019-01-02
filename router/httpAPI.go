package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/api/user"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

// SetUpHttpAPI set up http api for gin.Engine
func SetUpHttpAPI(r *gin.Engine) {

	common.ShouldPrintDebug = true
	httpAPI := r.Group("api/v1")
	{
		httpAPI.POST("/login", middleware.UserAuth.LoginHandler)
		httpAPI.POST("/register", user.Register)
		user.UserAPI(httpAPI)
	}
}
