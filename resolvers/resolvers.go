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
	r.crud = db.NewCRUD(mongoSession)
	return err
}

// CloseMongoDb close mongodb connection
func (r *RootResolver) CloseMongoDb() {
	r.crud.Close()
}

func bsonToAccount(b bson.M) models.Account {
	var account models.Account
	mapstructure.Decode(b, &account)
	account.ID = b["_id"].(bson.ObjectId)
	if b["hunter_id"] != nil {
		id := (b["hunter_id"]).(bson.ObjectId)
		account.HunterID = &id
	}
	if b["recruit_id"] != nil {
		id := (b["recruit_id"]).(bson.ObjectId)
		account.RecruitID = &id
	}
	return account
}
