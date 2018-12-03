package user

import (
	"../../common"
	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		c.JSON(500, "error")
		return
	}
	defer mdb.Close()

	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	col := mdb.DB("CUSmartFarm").C("users")
	col.Insert(gin.H{
		"username": username,
		"password": password,
		"email":    email,
	})
	c.JSON(200, "register ok")

}
