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

func TestLoginInvalid(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "login"
	queryFormat := `
		mutation {
			%s(%s) {
				refreshToken
				accessToken
			}
		}	  
	`
	acc := accounts[0]

	// prepare invalid input
	input := []string{
		fmt.Sprintf(`
			# case 1 invalid email
			email: "trash", password: "%s"
		`, acc.Password),
		fmt.Sprintf(`
			# case 2 invalid password
			email: "%s", password: "trash"
		`, acc.Email),
		`
			# case 3 invalid email and password
			email: "trash", password: "trash"
		`,
	}

	for i, in := range input {
		// request and respond
		query := fmt.Sprintf(queryFormat, method, in)
		response, err := gqlRequestAndRespond(handler, query, nil)

		//process response
		assert := assert.New(t)
		if err != nil {
			assert.Fail("Failed to process response:", err)
		}
		//make assertions
		assert.Contains(response, "errors", fmt.Sprintf("Case [%v]: %s", i+1, msgNoError))
	}
}
