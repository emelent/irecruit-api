package database

import (
	"os"

	config "../config"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateMongoURL creates a mongo url
func CreateMongoURL(user, pass, host, port string) string {
	url := host
	if user != "" && pass != "" {
		url = user + ":" + pass + "@" + host
	}

	if port != "" {
		url += ":" + port
	}
	return "mongodb://" + url
}

//NewCRUD creates a new CRUD type
func NewCRUD(session *mgo.Session) *CRUD {
	if session != nil {
		clone := session.Copy()
		defer clone.Close()

		ensureIndexes(clone)
	}
	crud := &CRUD{}
	crud.Session = session
	crud.TempStorage = make(map[string][]bson.M)
	return crud
}

var collectionIndexes = map[string][]mgo.Index{
	config.AccountsCollection: []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	},
	config.TokenManagersCollection: []mgo.Index{
		{
			Key:    []string{"account_id"},
			Unique: true,
		},
		{
			Key:    []string{"refresh_token"},
			Unique: true,
		},
	},
	config.IndustriesCollection: []mgo.Index{
		{
			Key:    []string{"name"},
			Unique: true,
		},
	},
	config.DocumentsCollection: []mgo.Index{
		{
			Key:    []string{"url"},
			Unique: true,
		},
	},
}

func ensureIndexes(session *mgo.Session) {
	for name, indexes := range collectionIndexes {
		c := session.DB(os.Getenv("DB_NAME")).C(name)
		for _, index := range indexes {
			c.EnsureIndex(index)
		}
	}
}
