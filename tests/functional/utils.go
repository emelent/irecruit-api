package functionaltests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	db "../../database"
	moc "../../mocks"
	route "../../routing"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------------
// Functional test general helper functions
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

func assertGqlData(method string, response map[string]interface{}, assert *assert.Assertions) map[string]interface{} {
	dataPortion, dOk := response["data"].(map[string]interface{})

	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	return dataPortion
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
