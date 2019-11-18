package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/storage"
)

// Init : start worker
func Init() {
	fmt.Println("[Worker] Starting worker")
	fmt.Println("[Worker] Loading devices from DB")
	mdb, err := common.Mongo()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer mdb.Close()

	it := mdb.DB("CUSmartFarm").C("devices").Find(nil).Iter()
	cnt := 0
	for cur := new(map[string]interface{}); it.Next(cur); cur = new(map[string]interface{}) {
		var dev device.Device
		dev.FromMap(*cur)
		storage.Devices[dev.ID] = &dev
		cnt++
	}

	fmt.Printf("[Worker] Done loading %d devices\n", cnt)
	for _, dev := range storage.Devices {
		fmt.Printf("%v\n", dev)
	}
}

// Work loop fuunction for worker
func Work() {
	for true {
		for _, dev := range storage.Devices {
			fmt.Printf("WORKER: Device %s has value of %#v", dev.ID, dev.LastSensorValues)
			dev.BroadCast("1.0", true)
			common.Println("Send message to", dev.ID)
		}
		time.Sleep(300 * time.Second)
	}
}

// return 60*hour + min
func minutes(hour, min int) int {
	return 60*hour + min
}
