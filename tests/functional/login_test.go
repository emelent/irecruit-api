package functionaltests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValid(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "login"
	query := fmt.Sprintf(`
		mutation {
			%s(email: "%s", password: "%s") {
				refreshToken
				accessToken
			}
		}	  
	`, method, accounts[0].Email, accounts[0].Password)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultTokens, rOk := dataPortion[method].(map[string]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Contains(resultTokens, "accessToken", msgMissingResponseData)
	assert.Contains(resultTokens, "refreshToken", msgMissingResponseData)
}
