package database

import (
	"fmt"

	config "../config"
	mgo "gopkg.in/mgo.v2"
)

//CRUD is a db abstraction layer used to perforom testing
//as well as interact with the mgo
type CRUD struct {
	Session     *mgo.Session
	CopySession *mgo.Session
	TempStorage map[string][]interface{}
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
		// !! assumes values is an array.. is it?
		db.TempStorage[collection] = append(db.TempStorage[collection], values...)
		fmt.Println(len(db.TempStorage[collection]), collection, "(s) added.")
		return nil
	}

	db.InitCopy()
	err := db.CopySession.DB(config.DbName).C(collection).Insert(values...)
	return err
}

//FindAll  finds all matching db entries
func (db *CRUD) FindAll(collection string, query interface{}) ([]interface{}, error) {
	//mock
	if db.Session == nil {
		return db.TempStorage[collection], nil
	}

	db.InitCopy()
	var results []interface{}
	err := db.CopySession.DB(config.DbName).C(collection).Find(query).All(&results)
	return results, err
}

//FindOne finds a db entry
func (db *CRUD) FindOne(collection string, query interface{}) (interface{}, error) {
	//mock
	if db.Session == nil {
		return db.TempStorage[collection], nil
	}

	db.InitCopy()
	var result interface{}
	err := db.CopySession.DB(config.DbName).C(collection).Find(query).One(&result)
	return result, err
}

//UpdateID updates entry by id
func (db *CRUD) UpdateID(collection string, id, value interface{}) error {
	//mock
	if db.Session == nil {
		//Not implemented because I don't know how to do it yet.
		//This needs me to check type of entry then make a temp entry,
		//check if it has an 'id' property or assume it has one and write
		//the necessary code for checks and fallbacks, and then finally update
		//the entry in the slice, which I think is a bit much for a mock,
		//and haven't come across a need for this to test my handlers.

		// items := db.TempStorage[collection]
		// for i, v := range items {
		// 	if v.id == string(id) {
		// 		temp := append(items[:n], append([]interface {
		// 			v
		// 		}, items[n+1:]...)...)
		// 	}
		// 	db.TempStorage[collection] = temp
		// }
		return nil
	}

	db.InitCopy()
	return db.CopySession.DB(config.DbName).C(collection).UpdateId(id, value)
}

//DeleteID deletes a db entry by id
func (db *CRUD) DeleteID(collection string, id interface{}) error {
	//mock
	if db.Session == nil {
		//This is not implemented for the same reason UpdateID is
		//not implemented.
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
