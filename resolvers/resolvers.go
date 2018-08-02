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
func bsonToAccount(b bson.M) models.Account {
	var account models.Account
	return account
}
