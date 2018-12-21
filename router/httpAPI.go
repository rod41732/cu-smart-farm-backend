package router

import (
	"../api/middleware"
	"../api/user"
	"../common"
	"github.com/gin-gonic/gin"
)

func SetUpHttpAPI(r *gin.Engine) {

	common.ShouldPrintDebug = true
	httpAPI := r.Group("api/v1")
	// define short name
	userAuth := middleware.UserAuth.MiddlewareFunc()
	// ownerCheck := middleware.OwnerCheck
	{
		httpAPI.POST("/login", middleware.UserAuth.LoginHandler)
		httpAPI.POST("/register", user.Register)
		httpAPI.POST("/ws", userAuth /* WS*/)
	}
}
