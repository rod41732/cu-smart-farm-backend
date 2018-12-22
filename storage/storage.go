package storage

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/model"
)

var mappedUserObject = make(map[string]model.User, 0)

// SetUserStateInfo : Map username into model.User
func SetUserStateInfo(username string, user model.User) {
	fmt.Printf("added user: %s\n", username)
	mappedUserObject[username] = user
}

// GetUserStateInfo get model.User corresponding to username
func GetUserStateInfo(username string) model.User {
	val, ok := mappedUserObject[username]
	fmt.Printf("get user: %s is ok=%v\n", username, ok)
	if !ok {
		val = nil
	}
	return val
}
