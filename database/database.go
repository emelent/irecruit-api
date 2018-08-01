package database

import mgo "gopkg.in/mgo.v2"

//NewCRUD creates a new CRUD type
func NewCRUD(session *mgo.Session) *CRUD {
	crud := &CRUD{}
	crud.Session = session
	crud.TempStorage = make(map[string][]interface{})
	return crud
}
