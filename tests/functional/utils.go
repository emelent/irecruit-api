package functionaltests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	config "../../config"
	db "../../database"
	moc "../../mocks"
	models "../../models"
	route "../../routing"
	utils "../../utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

// -------------------------------------------
// Helper functions
// -------------------------------------------

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
	// req.Header.Add("User-Agent", "go-tester")

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

// gqlRequestAndRespond query on handler and returns response
func gqlRequestAndRespond(handler http.Handler, query string, variables *map[string]interface{}) (map[string]interface{}, error) {
	req := createGqlRequest(query, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()

	response, err := getJSONResponse(res)
	return response, err
}

// createLoadedGqlHandler creates a new handler with a fully loaded crud
func createLoadedGqlHandler() http.Handler {
	crud := moc.NewLoadedCRUD()
	return createGqlHandler(crud)
}

// panicOnError panics on error
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// failOnError  fails test on error
func failOnError(assert *assert.Assertions, err error) {
	if err != nil {
		assert.Fail(err.Error())
	}
}

// login logs in as user specified by id setting the user-agent to value in ua
func login(crud *db.CRUD, id bson.ObjectId, ua string) (string, string) {
	accessToken, err := utils.CreateAccessToken(id.Hex(), ua)
	panicOnError(err)

	// get token mgr
	rawTokenMgr, err := crud.FindOne(config.TokenManagersCollection, &bson.M{"account_id": id})
	panicOnError(err)
	tokenManager := models.TransformTokenManager(rawTokenMgr)

	// TODO later when sys implemented,
	// store token in TokenMgr

	return accessToken, tokenManager.RefreshToken
}

// To be removed
func assertGqlData(method string, response map[string]interface{}, assert *assert.Assertions) map[string]interface{} {
	dataPortion, dOk := response["data"].(map[string]interface{})

	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	return dataPortion
}

// getSysUserAccount returns a system user account
func getSysUserAccount() models.Account {
	return moc.Accounts[len(moc.Accounts)-1]
}

// getNonSysUserAccount returns a non-system user account
func getNonSysUserAccount() models.Account {
	return moc.Accounts[0]
}

// getRecruitUserAccount returns a user account with a recruit profile
func getRecruitUserAccount() models.Account {
	return moc.Accounts[0]
}

// getNonRecruitUserAccount returns a user account without a recruit profile
func getNonRecruitUserAccount() models.Account {
	return moc.Accounts[2]
}

// getPlainUserAccount returns a non-sys user account without a hunter or recruit  profile
func getPlainUserAccount() models.Account {
	return moc.Accounts[len(moc.Accounts)-2]
}

// messages
const (
	msgUnexpectedError     = "Unexpected error in response."
	msgUnexpectedData      = "Unexpected data in response."
	msgInvalidResponse     = "Invalid response."
	msgMissingResponseData = "Missing response data."
	msgInvalidResponseType = "Invalid data response type."
	msgInvalidResultCount  = "Invalid number of results."
	msgInvalidResult       = "Invalid result."
	msgNoError             = "No outputted error."
)
