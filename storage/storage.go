package storage

import "github.com/rod41732/cu-smart-farm-backend/model"

var mappedUserObject map[string]model.User

// SetUserStateInfo : Map username into model.User
func SetUserStateInfo(username string, user model.User) {
	if mappedUserObject == nil {
		mappedUserObject = make(map[string]model.User, 0)
	}
	mappedUserObject[username] = user
}
