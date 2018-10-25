package router

import (
	"../common"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func MQTT() {

	common.CheckErr("connect to MQTT", common.ConnectToMQTT())

	client := &service.Client{}
	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("CUSmartFarm"), 1)

	subMsg.
}
