package utils

import (
	"log"
	"math/rand"

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

// PickRandomN create a random subset of length n from input slice
func PickRandomN(n int, input []interface{}) []interface{} {
	out := make([]interface{}, n)
	for i := range out {
		index := rand.Intn(len(input))
		out[i] = input[index]
		// remove used value from index
		input = append(input[:index], input[index+1:]...)

	}
	return out
}

// PickRandom pick a random value from input slice
func PickRandom(n int, input []interface{}) interface{} {
	return PickRandomN(1, input)[0]
}
