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
	assert.NotContains(response, "error", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Len(results, len(accounts), msgInvalidResultCount)
}

func TestRemoveAccount(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "removeAccount"
	query := fmt.Sprintf(`
		mutation{
			%s(
				id: "%s"
			)
		}
	`, method, accounts[0].ID.Hex())
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
	result, rOk := dataPortion[method].(string)

	//make assertions
	assert.NotContains(response, "error", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Equal(result, "Account successfully removed.", msgInvalidResult)
}
