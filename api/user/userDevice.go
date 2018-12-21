package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addDevice(c *gin.Context) {
	var match bson.M
	mdb, err := common.Mongo()
	defer mdb.Close()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	deviceID := c.PostForm("id")
	deviceSecret := c.PostForm("secret")
	common.Println(deviceID, deviceSecret)
	collection := mdb.DB("CUSmartFarm").C("devices")

	query := collection.Find(bson.M{
		"id":     deviceID,
		"secret": deviceSecret,
	})
	deviceQuery.One(&match)

	query.One(&match)

	if match == nil {
		return false, "Invalid device ID/ Secret"
	}
	if match["owner"] != nil {
		return false, "Device already owned"
	}

	collection.Update(bson.M{
		"username": username,
	}, bson.M{
		"$push": bson.M{
			"ownedDevices": deviceID,
		},
	})

	collection.Update(bson.M{
		"id":     deviceID,
		"secret": deviceSecret,
	}, bson.M{
		"$set": gin.H{
			"owner": username,
		},
	})

	c.JSON(200, gin.H{
		"msg": "added device",
	})
}

// HandleRemoveDevice handles removal of device (and check owner before doing so)
func HandleRemoveDevice(payload model.RemoveDeviceMessage, username string) (bool, string) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return false, "Something went wrong"
	}

	deviceID := c.PostForm("id")

	collection := mdb.DB("CUSmartFarm").C("devices")
	var match gin.H

	query := collection.Find(gin.H{
		"id": deviceID,
	})
	deviceQuery.One(&match)

	query.One(&match)
	if match == nil {
		return false, "Invalid DeviceID or Not Your Device"
	}

	var result interface{}
	removeDevice := mgo.Change{
		Update: bson.M{"$pull": bson.M{"ownedDevices": deviceID}},
	}

	var after interface{}

	collection.Find(gin.H{
		"username": user.(*middleware.User).Username,
	}).Apply(removeDevice, after)

	changeOwner := mgo.Change{
		Update: bson.M{"$set": bson.M{"owner": nil}},
	}
	query.Apply(changeOwner, after)
	c.JSON(200, "removed device")

	col.Find(bson.M{
		"username": username,
	}).Apply(removeDevice, result)
	deviceQuery.Apply(changeOwner, result)

	return true, "OK"
}
