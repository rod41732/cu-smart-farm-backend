package message

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/rod41732/cu-smart-farm-backend/common"
)

var provincesList = []string{
	"กระบี่", "กรุงเทพมหานคร", "กาญจนบุรี", "กาฬสินธุ์", "กำแพงเพชร", "ขอนแก่น", "จันทบุรี", "ฉะเชิงเทรา", "ชลบุรี", "ชัยนาท", "ชัยภูมิ", "ชุมพร", "ตรัง",
	"ตราด", "ตาก", "นครนายก", "นครปฐม", "นครพนม", "นครราชสีมา", "นครศรีธรรมราช", "นครสวรรค์", "นนทบุรี", "นราธิวาส", "น่าน", "บึงกาฬ", "บุรีรัมย์",
	"ปทุมธานี", "ประจวบคีรีขันธ์", "ปราจีนบุรี", "ปัตตานี", "พะเยา", "พังงา", "พัทลุง", "พิจิตร", "พิษณุโลก", "ภูเก็ต", "มหาสารคาม", "มุกดาหาร", "ยะลา",
	"ยโสธร", "ระนอง", "ระยอง", "ราชบุรี", "ร้อยเอ็ด", "ลพบุรี", "ลำปาง", "ลำพูน", "ศรีสะเกษ", "สกลนคร", "สงขลา", "สตูล", "สมุทรปราการ", "สมุทรสงคราม",
	"สมุทรสาคร", "สระบุรี", "สระแก้ว", "สิงห์บุรี", "สุพรรณบุรี", "สุราษฎร์ธานี", "สุรินทร์", "สุโขทัย", "หนองคาย", "หนองบัวลำภู", "อยุธยา", "อำนาจเจริญ",
	"อุดรธานี", "อุตรดิตถ์", "อุทัยธานี", "อุบลราชธานี", "อ่างทอง", "เชียงราย", "เชียงใหม่", "เพชรบุรี", "เพชรบูรณ์", "เลย", "แพร่", "แม่ฮ่องสอน"}

// EditProfileMessage payload for editing profile
type EditProfileMessage struct {
	Province string `json:"province"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// Validate validate payload format
func (message *EditProfileMessage) Validate() bool {
	match, _ := regexp.MatchString("[a-zA-Z0-9_]{3,}@[a-zA-Z0-9_]{3,}.[a-zA-Z0-9_]{3,}", message.Email)
	return match && common.StringInSlice(message.Province, provincesList)
}

//FromMap is "constructor" for converting map[string]interface{} to EditProfileMessage return error if can't convert
func (message *EditProfileMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	if message.Validate() {
		return nil
	} else {
		return errors.New("Validation Error")
	}
}

// ChangePasswordMessage payload for changing pasword
type ChangePasswordMessage struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Validate validate payload format
func (message *ChangePasswordMessage) Validate() bool {
	return len(message.OldPassword) >= 8 && len(message.NewPassword) >= 8
}

//FromMap is "constructor" for converting map[string]interface{} to EditProfileMessage return error if can't convert
func (message *ChangePasswordMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	if message.Validate() {
		return nil
	} else {
		return errors.New("Validation Error")
	}
}
