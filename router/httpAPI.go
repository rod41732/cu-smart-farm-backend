package router

import (
	"../api/device"
	"../api/user"
	"github.com/gin-gonic/gin"
)

func SetUpHttpAPI(r *gin.Engine) {

	httpAPI := r.Group("api/v1")
	device.DeviceControlAPI(httpAPI)
	user.SetUpUserAPI(httpAPI)
}
