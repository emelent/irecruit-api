package resolvers

import (
	db "../database"
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
