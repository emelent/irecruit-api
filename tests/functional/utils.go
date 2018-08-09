package functionaltests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// load data into db
	loadAccounts(crud)
	loadTokenManagers(crud)

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

		// validate before insertion
		if err := acc.OK(); err != nil {
			fmt.Printf("Mock accounts[%v] : %s", i, err.Error())
			break
		}

		// hash password before insertion
		acc.HashPassword()
		crud.Insert(accountsCollection, acc)
	}
}

func loadTokenManagers(crud *db.CRUD) {
	for i, mgr := range tokenManagers {
		mgr.AccountID = accounts[i].ID

		// validate before insertion
		if err := mgr.OK(); err != nil {
			fmt.Printf("Mock tokenManagers[%v] : %s", i, err.Error())
			break
		}

		crud.Insert(tokenManagersCollection, mgr)
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

func gqlRequestAndRespond(handler http.Handler, query string, variables *map[string]interface{}) (map[string]interface{}, error) {
	req := createGqlRequest(query, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()

	response, err := getJSONResponse(res)
	return response, err
}

// ================
// DATA
// ================
//

// collections
const (
	accountsCollection      = "accounts"
	tokenManagersCollection = "token_managers"
)

// messages
const (
	msgUnexpectedError     = "Unexpected error in response."
	msgInvalidResponse     = "Invalid response."
	msgMissingResponseData = "Missing response data."
	msgInvalidResponseType = "Invalid data response type."
	msgInvalidResultCount  = "Invalid number of results."
	msgInvalidResult       = "Invalid result."
	msgNoError             = "No outputted error."
)

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

// 6  token managers
var tokenManagers = []models.TokenManager{
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
}

// 6 user accounts
var accounts = []models.Account{
	{
		ID:          bson.NewObjectId(),
		Email:       "mark@gmail.com",
		Password:    "password",
		Name:        "Mark",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jdoe@gmail.com",
		Password:    "password",
		Name:        "John",
		Surname:     "Doe",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "lisa@gmail.com",
		Password:    "password",
		Name:        "Lisa",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "erin@gmail.com",
		Password:    "password",
		Name:        "Erin",
		Surname:     "Lona",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jake@gmail.com",
		Password:    "password",
		Name:        "Jake",
		Surname:     "Tinder",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "thato@gmail.com",
		Password:    "password",
		Name:        "Thato",
		Surname:     "Mopani",
		AccessLevel: 9, // system admin
	},
}
