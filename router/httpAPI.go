package router

import (
	"../api/user"
	"../common"
	"github.com/gin-gonic/gin"
)

func SetUpHttpAPI(r *gin.Engine) {

	common.ShouldPrintDebug = true
	httpAPI := r.Group("api/v1")
	// device.DeviceControlAPI(httpAPI)
	user.SetUpUserAPI(httpAPI)
}
