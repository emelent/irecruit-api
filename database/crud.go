package database

import (
	config "../config"
	"github.com/fatih/structs"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//CRUD is a db abstraction layer used to perforom testing
//as well as interact with the mgo
type CRUD struct {
	Session     *mgo.Session
	CopySession *mgo.Session
	TempStorage map[string][]bson.M
}

//InitCopy initialises a copy session if one is not ready
func (db *CRUD) InitCopy() {
	if db.Session != nil && db.CopySession == nil {
		db.CopySession = db.Session.Copy()
	}
}

//Insert inserts into db
func (db *CRUD) Insert(collection string, values ...interface{}) error {
	//mock
	if db.Session == nil {
		// turn values into bson values
		bValues := make([]bson.M, 0)
		for _, v := range values {
			bValues = append(bValues, makeBson(v))
		}
		db.TempStorage[collection] = append(db.TempStorage[collection], bValues...)
		return nil
	}

	db.InitCopy()
	err := db.CopySession.DB(config.DbName).C(collection).Insert(values...)
	return err
}

//FindAll  finds all matching db entries
func (db *CRUD) FindAll(collection string, query *bson.M) ([]interface{}, error) {
	if db.Session == nil { // a.k.a, we're in the mock
		results := filter(db.TempStorage[collection], matchQuery(query))
		return results, nil
	}

	db.InitCopy()
	var results []interface{}
	err := db.CopySession.DB(config.DbName).C(collection).Find(query).All(&results)
	return results, err
}

//FindOne finds a db entry
func (db *CRUD) FindOne(collection string, query *bson.M) (interface{}, error) {
	if db.Session == nil { // in the mock
		result := filterFirst(db.TempStorage[collection], matchQuery(query))
		return result, nil
	}

	db.InitCopy()
	var result interface{}
	err := db.CopySession.DB(config.DbName).C(collection).Find(query).One(&result)
	return result, err
}

//UpdateID updates entry by id
func (db *CRUD) UpdateID(collection string, id bson.ObjectId, updates bson.M) error {
	if db.Session == nil { // mocking
		for i, r := range db.TempStorage[collection] {
			if r["_id"] == id {
				for k, v := range updates {
					r[k] = v
				}
				db.TempStorage[collection][i] = r
				break
			}
		}
		return nil
	}

	db.InitCopy()
	return db.CopySession.DB(config.DbName).C(collection).UpdateId(id, updates)
}

//DeleteID deletes a db entry by id
func (db *CRUD) DeleteID(collection string, id bson.ObjectId) error {
	//mock
	if db.Session == nil {
		c := db.TempStorage[collection]
		for i, v := range db.TempStorage[collection] {
			if v["_id"] == id {
				db.TempStorage[collection] = append(c[:i], c[i+1:]...)
				break
			}
		}

		return nil
	}

	db.InitCopy()
	return db.CopySession.DB(config.DbName).C(collection).RemoveId(id)
}

//Close closes both the copy and the original db session
func (db *CRUD) Close() {
	if db.Session != nil {
		db.CloseCopy()
		db.Close()
		db.Session = nil
	}
}

//CloseCopy closes copy db session
func (db *CRUD) CloseCopy() {
	if db.CopySession != nil {
		db.CopySession.Close()
		db.CopySession = nil
	}
}

func makeBson(in interface{}) bson.M {
	structs.DefaultTagName = "bson"
	t := structs.Map(in)
	return bson.M(t)
}

func filter(in []bson.M, fn func(bson.M) bool) []interface{} {
	results := make([]interface{}, 0)
	for _, v := range in {
		if fn(v) {
			results = append(results, v)
		}
	}
	return results
}

func filterFirst(in []bson.M, fn func(bson.M) bool) interface{} {
	for _, v := range in {
		if fn(v) {
			return v
		}
	}
	return nil
}
func matchQuery(query *bson.M) func(bson.M) bool {
	if query == nil {
		return func(m bson.M) bool {
			return true
		}
	}

	return func(m bson.M) bool {
		for k, v := range *query {
			if m[k] != v {
				return false
			}
		}
		return true
	}
}
