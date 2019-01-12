package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/rod41732/cu-smart-farm-backend/common"

	"github.com/rod41732/cu-smart-farm-backend/service/receiver"
	"github.com/rod41732/cu-smart-farm-backend/service/worker"
)

func main() {
	common.ShouldPrintDebug = true
	trigger := new(receiver.Trigger)
	rpc.Register(trigger)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	worker.Init()
	go worker.Work()
	http.Serve(l, nil)
}
