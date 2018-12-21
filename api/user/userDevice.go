package user

import (
	"../../common"
	"../middleware"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addDevice(c *gin.Context) {
	var match gin.H
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		c.JSON(500, gin.H{
			"msg": "Connection to database failed",
		})
		return
	}

	deviceID := c.PostForm("id")
	deviceSecret := c.PostForm("secret")
	common.Println(deviceID, deviceSecret)
	col := mdb.DB("CUSmartFarm").C("devices")

	que := col.Find(gin.H{
		"id":     deviceID,
		"secret": deviceSecret,
	})

	que.One(&match)

	if match == nil {
		c.JSON(404, gin.H{
			"msg": "device not found.",
		})
		return
	}

	user, _ := c.Get("username")
	username := user.(*middleware.User).Username

	if match["owner"] != nil && match["owner"] != username {
		c.JSON(403, gin.H{
			"msg": "device already owned",
		})
		common.Println(match["owner"].(string) + "/" + username)
		return
	} else if match["owner"] == username {
		c.JSON(200, gin.H{
			"msg": "already owned",
		})
		return
	}

	appendDevice := mgo.Change{
		Update: bson.M{
			"$push": bson.M{
				"ownedDevices": deviceID,
			},
		},
	}

	var after interface{}

	col.Find(gin.H{
		"username": username,
	}).Apply(appendDevice, after)

	changeOwner := mgo.Change{
		Update: gin.H{
			"$set": gin.H{
				"owner": username,
			},
		},
	}
	que.Apply(changeOwner, after)
	c.JSON(200, "added device")
}

func removeDevice(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		c.JSON(500, "error")
		return
	}

	deviceID := c.PostForm("id")

	col := mdb.DB("CUSmartFarm").C("devices")
	var match gin.H

	que := col.Find(gin.H{
		"id": deviceID,
	})

	que.One(&match)
	if match == nil {
		c.JSON(404, "userdevice: device not found")
		return
	}

	user, _ := c.Get("username")
	removeDevice := mgo.Change{
		Update: gin.H{
			"$pull": gin.H{
				"ownedDevices": deviceID,
			},
		},
	}

	var after interface{}

	col.Find(gin.H{
		"username": user.(*middleware.User).Username,
	}).Apply(removeDevice, after)

	changeOwner := mgo.Change{
		Update: gin.H{
			"$set": gin.H{
				"owner": nil,
			},
		},
	}
	que.Apply(changeOwner, after)
	c.JSON(200, "removed device")

}
