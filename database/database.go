package database

import (
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
	"accounts": []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	},
	"token_managers": []mgo.Index{
		{
			Key:    []string{"account_id"},
			Unique: true,
		},
		{
			Key:    []string{"refresh_token"},
			Unique: true,
		},
	},
	"industries": []mgo.Index{
		{
			Key:    []string{"name"},
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
