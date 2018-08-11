package resolvers

import (
	db "../database"
	models "../models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

// RootResolver contains functions that resolve graphql queries
type RootResolver struct {
	crud *db.CRUD
}

// Init initialises the crud system
func (r *RootResolver) Init(crud *db.CRUD) {
	if crud == nil {
		// create a mock CRUD instance if nil provided
		crud = db.NewCRUD(nil)
	}

	r.crud = crud
}

// transform interface to Account model
func transformAccount(in interface{}) models.Account {
	var account models.Account
	switch v := in.(type) {
	case bson.M:

		mapstructure.Decode(v, &account)
		account.ID = v["_id"].(bson.ObjectId)
		account.HunterID = (v["hunter_id"]).(bson.ObjectId)
		account.RecruitID = (v["recruit_id"]).(bson.ObjectId)
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
		tokenMgr.RefreshToken = v["refresh_token"].(string)
		tokenMgr.MaxTokens = v["max_tokens"].(int)
		tokenMgr.ID = v["_id"].(bson.ObjectId)
		tokenMgr.AccountID = v["account_id"].(bson.ObjectId)
		tokenMgr.Tokens = v["tokens"].([]string)
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
	case models.Recruit:
		recruit = v
	}

	return recruit
}
