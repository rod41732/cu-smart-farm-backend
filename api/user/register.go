package user

import (
	"fmt"
	"regexp"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

var provincesList = []string{
	"กระบี่", "กรุงเทพมหานคร", "กาญจนบุรี", "กาฬสินธุ์", "กำแพงเพชร", "ขอนแก่น", "จันทบุรี", "ฉะเชิงเทรา", "ชลบุรี", "ชัยนาท", "ชัยภูมิ", "ชุมพร", "ตรัง",
	"ตราด", "ตาก", "นครนายก", "นครปฐม", "นครพนม", "นครราชสีมา", "นครศรีธรรมราช", "นครสวรรค์", "นนทบุรี", "นราธิวาส", "น่าน", "บึงกาฬ", "บุรีรัมย์",
	"ปทุมธานี", "ประจวบคีรีขันธ์", "ปราจีนบุรี", "ปัตตานี", "พะเยา", "พังงา", "พัทลุง", "พิจิตร", "พิษณุโลก", "ภูเก็ต", "มหาสารคาม", "มุกดาหาร", "ยะลา",
	"ยโสธร", "ระนอง", "ระยอง", "ราชบุรี", "ร้อยเอ็ด", "ลพบุรี", "ลำปาง", "ลำพูน", "ศรีสะเกษ", "สกลนคร", "สงขลา", "สตูล", "สมุทรปราการ", "สมุทรสงคราม",
	"สมุทรสาคร", "สระบุรี", "สระแก้ว", "สิงห์บุรี", "สุพรรณบุรี", "สุราษฎร์ธานี", "สุรินทร์", "สุโขทัย", "หนองคาย", "หนองบัวลำภู", "อยุธยา", "อำนาจเจริญ",
	"อุดรธานี", "อุตรดิตถ์", "อุทัยธานี", "อุบลราชธานี", "อ่างทอง", "เชียงราย", "เชียงใหม่", "เพชรบุรี", "เพชรบูรณ์", "เลย", "แพร่", "แม่ฮ่องสอน"}

func validateNationalID(ID []byte) bool {
	if len(ID) != 13 {
		return false
	}
	sum := 0
	last := int(ID[12] - '0')
	for idx, char := range ID[0:12] {
		sum += (13 - idx) * (int(char - '0'))
	}
	return (11-(sum%11))%10 == last
}

// Register : API for user register
func Register(c *gin.Context) {
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At API/Register - Connecting to DB")
		c.JSON(500, "error")
		return
	}
	defer mdb.Close()

	username := c.PostForm("username")
	password := c.PostForm("password")
	province := c.PostForm("province")
	address := c.PostForm("address")
	nationalID := c.PostForm("nationalID")
	email := c.PostForm("email")
	cnt, err := mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": username,
	}).Count()
	if cnt != 0 {
		c.JSON(401, "username already exists")
		return
	}
	var ok = true
	var errmsg = "OK"
	// Validations, response first error
	if match, _ := regexp.MatchString("[a-zA-Z0-9_]{6,}", username); ok && !match {
		ok, errmsg = false, "Username error"
	}
	if ok && len(password) < 8 {
		ok, errmsg = false, "Password error"
	}
	if ok && !common.StringInSlice(province, provincesList) {
		ok, errmsg = false, "Province error"
	}
	if match, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, email); ok && !match {
		ok, errmsg = false, "Email error"
	}
	if ok && !validateNationalID([]byte(nationalID)) {
		ok, errmsg = false, "ID Error"
	}

	if ok {
		col := mdb.DB("CUSmartFarm").C("users")
		col.Insert(gin.H{
			"username":   username,
			"password":   common.SHA256(password),
			"province":   province,
			"address":    address,
			"nationalID": nationalID,
			"email":      email,
			"devices":    []string{},
		})
	}
	c.JSON(200, gin.H{
		"message": errmsg,
		"success": ok,
	})

}
