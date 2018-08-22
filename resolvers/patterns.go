package resolvers

import (
	"log"

	config "../config"
	db "../database"
	er "../errors"
	models "../models"
	"gopkg.in/mgo.v2/bson"
)

// ResolveRemoveByID is a generic resolver for removeByID methods
func ResolveRemoveByID(crud *db.CRUD, collection, name, id string) (*string, error) {
	defer crud.CloseCopy()

	// check that the ID is valid
	if !bson.IsObjectIdHex(id) {
		return nil, er.InvalidField("id")
	}

	// attempt to remove document
	if err := crud.DeleteID(collection, bson.ObjectIdHex(id)); err != nil {
		return nil, er.Generic()
	}
	result := name + " successfully removed."
	return &result, nil
}

// GenericUpdateByID performs a generic update and returns the new result
func GenericUpdateByID(crud *db.CRUD, collection string, id bson.ObjectId, updates bson.M) (interface{}, error) {
	defer crud.CloseCopy()

	// attempt update
	if err := crud.UpdateID(collection, id, updates); err != nil {
		return nil, er.Generic()
	}
	result, err := crud.FindID(collection, id)
	if err != nil {
		return nil, er.Generic()
	}
	return result, nil
}

// ResolveRemoveAccount is a generic resolver for removing an account by ID along
// along with the corresponding TokenManager
func ResolveRemoveAccount(crud *db.CRUD, id bson.ObjectId) (*string, error) {
	defer crud.CloseCopy()

	// delete the account
	err := crud.DeleteID(config.AccountsCollection, id)
	if err != nil {
		log.Println("Failed to delete Account =>", err)
		return nil, er.Generic()
	}

	// find the account's token manager
	rawTokenMgr, err := crud.FindOne(config.TokenManagersCollection, bson.M{"account_id": id})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.Generic()
	}

	// delete the account's token manager
	tokenMgr := models.TransformTokenManager(rawTokenMgr)
	err = crud.DeleteID(config.TokenManagersCollection, tokenMgr.ID)
	if err != nil {
		log.Println("Failed to delete TokenManager =>", err)
		return nil, er.Generic()
	}

	msg := "Account successfully removed."
	return &msg, nil
}
