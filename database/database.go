package database

import (
	config "../config"
	mgo "gopkg.in/mgo.v2"
)

//NewCRUD creates a new CRUD type
func NewCRUD(session *mgo.Session) *CRUD {
	clone := session.Copy()
	defer clone.Close()

	ensureIndexes(clone)
	crud := &CRUD{}
	crud.Session = session
	crud.TempStorage = make(map[string][]interface{})
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
			Key:    []string{"tokens.signature"},
			Unique: true,
		},
	},
}

func ensureIndexes(session *mgo.Session) {
	for name, indexes := range collectionIndexes {
		c := session.DB(config.DbName).C(name)
		for _, index := range indexes {
			c.EnsureIndex(index)
		}
	}
}
