package user

import (
	"../../common"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func register(c *gin.Context) {
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
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if common.PrintError(err) {
		c.JSON(500, "Error registering")
		return
	}
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
