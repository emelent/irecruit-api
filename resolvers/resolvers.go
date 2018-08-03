package resolvers

import (
	config "../config"
	db "../database"
	models "../models"
	"github.com/mitchellh/mapstructure"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RootResolver contains functions that resolve graphql queries
type RootResolver struct {
	crud *db.CRUD
}

// OpenMongoDb opens mongodb connection
func (r *RootResolver) OpenMongoDb() error {
	mongoSession, err := mgo.Dial(config.DbHost)
	if err == nil {
		r.crud = db.NewCRUD(mongoSession)
	} else {
		r.crud = db.NewCRUD(mongoSession)
	}

	return err
}

// CloseMongoDb close mongodb connection
func (r *RootResolver) CloseMongoDb() {
	r.crud.Close()
}

func transformAccount(in interface{}) models.Account {
	var account models.Account
	switch v := in.(type) {
	case bson.M:

		mapstructure.Decode(v, &account)
		account.ID = v["_id"].(bson.ObjectId)
		if v["hunter_id"] != nil {
			id := (v["hunter_id"]).(bson.ObjectId)
			account.HunterID = &id
		}
		if v["recruit_id"] != nil {
			id := (v["recruit_id"]).(bson.ObjectId)
			account.RecruitID = &id
		}

	case models.Account:
		account = v
	}

	return account
}

func transformTokenManager(in interface{}) models.TokenManager {
	var tokenMgr models.TokenManager
	switch v := in.(type) {
	case bson.M:
		tokenMgr.RefreshToken = v["refresh_token"].(string)
		tokenMgr.MaxTokens = v["max_tokens"].(int)
		tokenMgr.ID = v["_id"].(bson.ObjectId)
		tokenMgr.AccountID = v["account_id"].(bson.ObjectId)
		tokens := make([]string, 0)
		rawTokens := v["tokens"].([]interface{})
		for _, t := range rawTokens {
			tokens = append(tokens, t.(string))
		}
		tokenMgr.Tokens = tokens
	case models.TokenManager:
		tokenMgr = v
	}

	return tokenMgr
}
