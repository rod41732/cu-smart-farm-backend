package user

import (
	"../../common"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		c.JSON(500, "error")
		return
	}
	defer mdb.Close()

	username := c.PostForm("username")
	password := common.SHA256(c.PostForm("password"))
	province := c.PostForm("province")
	address := c.PostForm("address")
	nationalID := c.PostForm("nationalID")
	email := c.PostForm("email")

	col := mdb.DB("CUSmartFarm").C("users")
	col.Insert(gin.H{
		"username":   username,
		"password":   password,
		"province":   province,
		"address":    address,
		"nationalID": nationalID,
		"email":      email,
	})
	c.JSON(200, gin.H{
		"status": "OK",
	})

}
