package user

import (
	"../middleware"
	"github.com/gin-gonic/gin"
)

func SetUpUserAPI(r *gin.RouterGroup) {
	userAPI := r.Group("/user")
	{
		userAPI.POST("/login", middleware.UserAuth.LoginHandler)
	}
}
