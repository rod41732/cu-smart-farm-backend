package user

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/rod41732/cu-smart-farm-backend/model/message"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

// GetDevList : get user device list
func (user *RealUser) GetDevList() {
	resp, err := json.Marshal(gin.H{
		"t":       "response",
		"e":       "getDevList",
		"payload": user.Devices,
	})
	if common.PrintError(err) {
		fmt.Println("  At User::GetDevList -- json.Marshal")
	}
	user.conn.WriteMessage(1, resp)
}

// EditProfile edits user's profile in DB
func (user *RealUser) EditProfile(payload map[string]interface{}) (bool, string) {
	var message message.EditProfileMessage
	if message.FromMap(payload) != nil {
		return false, "Bad Payload"
	}

	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At User::EditProfile -- Connecting to DB")
		return false, "Can't connect to DB"
	}
	defer mdb.Close()
	var tmp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": user.Username,
	}).Apply(mgo.Change{
		Update: bson.M{
			"$set": message,
		},
	}, &tmp)
	if common.PrintError(err) {
		fmt.Println("  At User::EditProfile -- Updating Info")
		return false, "User modify error"
	}
	user.Province = message.Province
	user.Address = message.Address
	user.Email = message.Email
	return true, "OK"
}

// ChangePassword changes user's password, (require confirming old password)
func (user *RealUser) ChangePassword(payload map[string]interface{}) (bool, string) {
	var message message.ChangePasswordMessage
	if message.FromMap(payload) != nil {
		return false, "Bad Payload"
	}
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At User::ChangePassword -- Connecting to DB")
		return false, "Can't connect to DB"
	}
	defer mdb.Close()
	var tmp map[string]interface{}
	_, err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": user.Username,
		"password": common.SHA256(message.OldPassword),
	}).Apply(mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"password": common.SHA256(message.NewPassword),
			},
		},
	}, &tmp)
	if common.PrintError(err) {
		fmt.Println("  At User::ChangePassword -- Updating Password")
		return false, "User modify error"
	}
	return true, "OK"
}
