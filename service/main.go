package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/router"

	"github.com/rod41732/cu-smart-farm-backend/mqtt"

	"github.com/rod41732/cu-smart-farm-backend/common"

	"github.com/rod41732/cu-smart-farm-backend/service/receiver"
	"github.com/rod41732/cu-smart-farm-backend/service/worker"
)

func main() {
	common.ShouldPrintDebug = true
	trigger := new(receiver.Trigger)
	rpc.Register(trigger)
	rpc.HandleHTTP()

	common.InitializeKeyPair()
	middleware.Initialize()

	router.InitMQTT()
	go mqtt.MQTT()
	time.Sleep(2 * time.Second) // wait until MQTT connect
	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	worker.Init()
	go worker.Work()
	http.Serve(l, nil)
}
