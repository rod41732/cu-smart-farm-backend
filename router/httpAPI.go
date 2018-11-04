package router

import (
	"../api/device"
	"github.com/gin-gonic/gin"
)

func SetUpHttpAPI(r *gin.Engine) {

	httpAPI := r.Group("api/v1")

	device.DeviceControlAPI(httpAPI)

}
