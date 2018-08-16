package utils

import (
	"log"

	models "../models"
	"gopkg.in/mgo.v2/bson"
)

// IsSysAccount checks if given account is a system account
func IsSysAccount(account *models.Account) bool {
	if account == nil {
		return false
	}
	log.Println("AccessLevel =>", account.AccessLevel)
	log.Println("isSys =>", account.AccessLevel > 5)
	return account.AccessLevel > 5
}

// IsNullID checks if given id references a "null" id
func IsNullID(id bson.ObjectId) bool {
	return id == models.NullObjectID
}
