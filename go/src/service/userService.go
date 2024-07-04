package service

import (
	"shorturl/go/src/db"
	"shorturl/go/src/modle"
)

func GetUserInfo(username string) (user modle.User, err error) {
	query := "SELECT * FROM `users` WHERE `username` = ?"
	result := db.Db.Raw(query, username).First(&user)
	return user, result.Error
}
