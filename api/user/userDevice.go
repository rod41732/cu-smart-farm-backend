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
	if common.PrintError(err) {
		c.JSON(500, gin.H{
			"msg": "Connection to database failed",
		})
		return
	}

	deviceID := c.PostForm("id")
	deviceSecret := c.PostForm("secret")
	common.Println(deviceID, deviceSecret)
	collection := mdb.DB("CUSmartFarm").C("devices")

	query := collection.Find(bson.M{
		"id":     deviceID,
		"secret": deviceSecret,
	})

	query.One(&match)

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

func removeDevice(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		c.JSON(500, "error")
		return
	}

	deviceID := c.PostForm("id")

	collection := mdb.DB("CUSmartFarm").C("devices")
	var match gin.H

	query := collection.Find(gin.H{
		"id": deviceID,
	})

	query.One(&match)
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

	collection.Find(gin.H{
		"username": user.(*middleware.User).Username,
	}).Apply(removeDevice, after)

	changeOwner := mgo.Change{
		Update: gin.H{
			"$set": gin.H{
				"owner": nil,
			},
		},
	}
	query.Apply(changeOwner, after)
	c.JSON(200, "removed device")

}
