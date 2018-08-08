package functionaltests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountList(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "accounts"
	query := fmt.Sprintf(`query{%s{id,name,surname,email}}`, method)
	req := createGqlRequest(query, nil)

	//make request
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()

	//process response
	assert := assert.New(t)
	response, err := getJSONResponse(res)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	results, rOk := dataPortion[method].([]interface{})

	//make assertions
	assert.NotContains(response, "error", "Unexpected error in response.")
	assert.Contains(response, "data", "Invalid response.")
	assert.Contains(response["data"], "accounts", "Missing response data.")
	assert.True(dOk, "Invalid data response type.")
	assert.True(rOk, "Invalid data[\"account\"] response type.")
	assert.Len(results, len(accounts), "Invalid number of results.")
}
