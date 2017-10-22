package auth

import (
	"donategold.me/components/users"
	"donategold.me/db"
)

func GetProviderByTypeAndUid(provider, uid string) (users.Provider, bool) {
	var p users.Provider
	exist := false
	db.DB.Where("type_provider = ? and uid = ?", provider, uid).First(&p)
	if p.ID != 0 {
		exist = true
	}
	return p, exist
}

func GetUserByProvider(p users.Provider) users.User {
	var user users.User
	user, exist := users.GetUserById(p.UserID)
	if !exist {
		user = users.User{Username: p.Uid}
	}
	return user
}
