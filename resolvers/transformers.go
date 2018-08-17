package resolvers

import (
	models "../models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

// transform interface to Account model
func transformAccount(in interface{}) models.Account {
	var account models.Account
	switch v := in.(type) {
	case bson.M:
		account.ID = v["_id"].(bson.ObjectId)
		account.Email = v["email"].(string)
		account.Name = v["name"].(string)
		account.Password = v["password"].(string)
		account.AccessLevel = v["access_level"].(int)
		account.HunterID = v["hunter_id"].(bson.ObjectId)
		account.RecruitID = v["recruit_id"].(bson.ObjectId)

	case map[string]interface{}:
		mapstructure.Decode(v, &account)

	case models.Account:
		account = v
	}

	return account
}

// transform  interface into TokenManager model
func transformTokenManager(in interface{}) models.TokenManager {
	var tokenMgr models.TokenManager
	switch v := in.(type) {
	case bson.M:
		tokenMgr.ID = v["_id"].(bson.ObjectId)
		tokenMgr.RefreshToken = v["refresh_token"].(string)
		tokenMgr.MaxTokens = v["max_tokens"].(int)
		tokenMgr.AccountID = v["account_id"].(bson.ObjectId)
		tokenMgr.Tokens = v["tokens"].([]string)

	case map[string]interface{}:
		mapstructure.Decode(v, &tokenMgr)

	case models.TokenManager:
		tokenMgr = v
	}

	return tokenMgr
}

// transform interface into Recruit model
func transformRecruit(in interface{}) models.Recruit {
	var recruit models.Recruit
	switch v := in.(type) {
	case bson.M:
		mapstructure.Decode(v, &recruit)
		recruit.ID = v["_id"].(bson.ObjectId)

	case map[string]interface{}:
		mapstructure.Decode(v, &recruit)

	case models.Recruit:
		recruit = v
	}

	return recruit
}

// transform interface into Industry model
func transformIndustry(in interface{}) models.Industry {
	var industry models.Industry
	switch v := in.(type) {
	case bson.M:
		mapstructure.Decode(v, &industry)

	case map[string]interface{}:
		mapstructure.Decode(v, &industry)

	case models.Industry:
		industry = v
	}

	return industry
}
