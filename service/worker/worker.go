package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/service/storage"
)

// Init : start worker
func Init() {
	fmt.Println("[Worker] Starting worker")
	fmt.Println("[Worker] Loading devices from DB")
	mdb, err := common.Mongo()
	if err != nil {
		log.Fatalln(err.Error())
	}

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
		toDevice := make(map[string]device.RelayState)
		for _, dev := range storage.Devices {
			for rID, state := range dev.RelayStates {
				var sched device.ScheduleDetail
				detailMap, ok := state.Detail.(map[string]interface{})
				if ok {
					err := sched.FromMap(detailMap)
					if err == nil {
						detialStr := "off"

						t := time.Now()
						now := minutes(t.Hour(), t.Minute())
						for _, entry := range sched.Schedules {
							if minutes(entry.StartHour, entry.StartMin) <= now && now <= minutes(entry.EndHour, entry.EndMin) {
								detialStr = "on"
								break
							}
						}
						toDevice[rID] = device.RelayState{
							Mode:   "manual",
							Detail: detialStr,
						}
					}
				}
				// TODO: repeat=false logic
			}
			if len(toDevice) > 0 {
				str, _ := json.Marshal(toDevice)
				common.Printf("[Worker] >>> send message to %s\n", dev.ID)
				dev.SendMsg(str)
			}
		}
		time.Sleep(60 * time.Second)
	}
}

// return 60*hour + min
func minutes(hour, min int) int {
	return 60*hour + min
}
