package functionaltests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gopkg.in/mgo.v2/bson"

	db "../../database"
	models "../../models"
	route "../../routing"
)

// Functional test general helper functions

// loadedCrud returns a crud object loaded with all the data
func loadedCrud() *db.CRUD {
	crud := db.NewCRUD(nil)

	// load accounts into db
	loadAccounts(crud)

	return crud
}

// loadAccounts loads all account data
func loadAccounts(crud *db.CRUD) {
	numHunters := len(hunterIDs)
	numRecruits := len(recruitIDs)

	for i, acc := range accounts {
		if i < numHunters { // first n users are hunters
			acc.HunterID = &(hunterIDs[i])
		} else if i < numHunters+numRecruits { // next m users are recruits
			acc.RecruitID = &(recruitIDs[(i - numHunters)])
		}

		crud.Insert(accountsCollection, acc)
	}
}

// createGqlHandler creates a graphql handler
func createGqlHandler(crud *db.CRUD) http.Handler {
	return route.NewGqlHandler(crud)
}

// createGqlRequest creates a graphql http.Request object
func createGqlRequest(query string, variables *map[string]interface{}) *http.Request {
	data := map[string]interface{}{
		"query": query,
	}
	if variables != nil {
		data["variables"] = *variables
	}

	postData, _ := json.Marshal(data)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(postData))
	req.Header.Add("Content-Type", "application/json")
	return req
}

// getJSONResponse unmarshals a map from a JSON type response
func getJSONResponse(res *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ================
// DATA
// ================
//

const accountsCollection = "accounts"

// 3 hunter IDs
var hunterIDs = []bson.ObjectId{
	bson.NewObjectId(),
	bson.NewObjectId(),
	bson.NewObjectId(),
}

// 2 recruit IDs
var recruitIDs = []bson.ObjectId{
	bson.NewObjectId(),
	bson.NewObjectId(),
}

// 6 user accounts
var accounts = []models.Account{
	{
		ID:          bson.NewObjectId(),
		Email:       "mark@gmail.com",
		Password:    "123",
		Name:        "Mark",
		Surname:     "Smith",
		AccessLevel: 9,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jdoe@gmail.com",
		Password:    "123",
		Name:        "John",
		Surname:     "Doe",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "lisa@gmail.com",
		Password:    "123",
		Name:        "Lisa",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "erin@gmail.com",
		Password:    "123",
		Name:        "Erin",
		Surname:     "Lona",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jake@gmail.com",
		Password:    "123",
		Name:        "Jake",
		Surname:     "Tinder",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "thato@gmail.com",
		Password:    "123",
		Name:        "Thato",
		Surname:     "Mopani",
		AccessLevel: 0,
	},
}
