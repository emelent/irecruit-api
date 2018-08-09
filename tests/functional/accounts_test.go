package functionaltests

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

func TestAccountList(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "accounts"
	query := fmt.Sprintf(`query{%s{id,name,surname,email}}`, method)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultAccounts, rOk := dataPortion[method].([]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Len(resultAccounts, len(accounts), msgInvalidResultCount)
}

func TestRemoveAccountValid(t *testing.T) {
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

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	result, rOk := dataPortion[method].(string)

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Equal(result, "Account successfully removed.", msgInvalidResult)
}

func TestRemoveAccountInvalid(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "removeAccount"
	queryFormat := `
		mutation{
			%s(
				id: "%s"
			)
		}
	`

	input := []string{
		"123",
		bson.NewObjectId().Hex(),
	}
	//make request
	for _, in := range input {
		query := fmt.Sprintf(queryFormat, method, in)
		response, err := gqlRequestAndRespond(handler, query, nil)

		//process response
		assert := assert.New(t)
		if err != nil {
			assert.Fail("Failed to process response:", err)
		}

		//make assertions
		assert.Contains(response, "errors", msgNoError)
	}

}

func TestCreateAccountValid(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "createAccount"
	queryFormat := `
		mutation{
			%s(%s){
				refreshToken
				accessToken
			}
		}
	`

	input := `
		info: {
			email: "test@gmail.com",
			password:"test",
			name: "Test",
			surname:"User"
		}
	`
	query := fmt.Sprintf(queryFormat, method, input)
	//make request
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

func TestCreateAccountInvalid(t *testing.T) {
	//prepare handler
	crud := loadedCrud()
	handler := createGqlHandler(crud)

	//prepare request
	method := "createAccount"
	queryFormat := `
		mutation{
			%s(%s){
				refreshToken
				accessToken
			}
		}
	`

	input := []string{
		`
			# missing field, (email)
			info: {
				password:"test",
				name: "Test",
				surname:"User"
			}
		`,
		`
			# invalid data type
			info: {
				email: 14,
				password:"test",
				name: "Test",
				surname:"User"
			}	
		`,
		`
			# invalid data content
			info: {
				email: "marshia",
				password:"test",
				name: "Test",
				surname:"User"
			}	
		`,

		// can't test duplicate entries because mock doesn't have that infrastructure
		// error is on the db layer
		// fmt.Sprintf(`
		// 	# duplicate email
		// 	info: {
		// 		email: "%s",
		// 		password:"test",
		// 		name: "Test",
		// 		surname:"User"
		// 	}
		// `, accounts[0].Email),
	}
	for i, in := range input {
		query := fmt.Sprintf(queryFormat, method, in)
		//make request
		response, err := gqlRequestAndRespond(handler, query, nil)

		//process response
		assert := assert.New(t)
		if err != nil {
			assert.Fail("Failed to process response:", err)
		}

		//make assertions
		assert.Contains(response, "errors", fmt.Sprintf("Case [%v]: %s", i, msgUnexpectedError))
	}

}
