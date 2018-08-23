package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestRandomQuestionsValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	//prepare request
	method := "randomQuestions"
	query := fmt.Sprintf(`
		query{
			%s(industry_id: "%s"){
				id,
				question,
				industry_id
			}
		}`, method, moc.Industries[0].ID.Hex())

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultQuestions, rOk := dataPortion[method].([]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	if rOk {
		assert.Len(resultQuestions, 2, msgInvalidResultCount)
	}
}

func TestRandomQuestionsInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "randomQuestions"
	queryFormat := `
		query{
			%s(%s){
				id,
				question,
				industry_id
			}
		}
	`

	// invalid inputs
	input := []string{
		`
			# case 1 no industry_id
		`,
		`
			# case 2 invalid industry_id
			industry_id: "43"
		`,
	}

	for i, in := range input {
		query := fmt.Sprintf(queryFormat, method, in)
		// request and respond
		response, err := gqlRequestAndRespond(handler, query, nil)
		//process response
		assert := assert.New(t)
		if err != nil {
			assert.Fail("Failed to process response:", err)
		}

		assert.Contains(response, "errors", fmt.Sprintf("Case [%v]: %s", i+1, msgNoError))
	}
}

func TestLoginValid(t *testing.T) {
	//prepare handler
	crud := moc.NewLoadedCRUD()
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
	`, method, moc.Accounts[0].Email, moc.DefaultPassword)

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
	crud := moc.NewLoadedCRUD()
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
	acc := moc.Accounts[0]

	// prepare invalid input
	input := []string{
		fmt.Sprintf(`
			# case 1 invalid email
			email: "trash", password: "%s"
		`, moc.DefaultPassword),
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

func TestCreateAccountValid(t *testing.T) {
	//prepare handler
	crud := moc.NewLoadedCRUD()
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
			password:"password"
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
	crud := moc.NewLoadedCRUD()
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

	// invalid inputs
	input := []string{
		`
			# case 1 missing field, (email)
			info: {
				password:"password"
				name: "Test",
				surname:"User"
			}
		`,
		`
			# case 2 invalid data type
			info: {
				email: 14,
				password:"password"
				name: "Test",
				surname:"User"
			}	
		`,
		`
			# case 3 invalid email
			info: {
				email: "marshia",
				password:"password"
				name: "Test",
				surname:"User"
			}	
		`,
		`
			# case 4 short password
			info: {
				email: "marshia@gmail.com",
				password:"123"
				name: "Test",
				surname:"User"
			}	
		`,
		`
			# case 5 short name
			info: {
				email: "marshia@gmail.com",
				password:"password"
				name: "T",
				surname:"User"
			}	
		`,
		`
			# case 6 short surname
			info: {
				email: "marshia@gmail.com",
				password:"password"
				name: "Test",
				surname:"u"
			}	
		`,

		// can't test duplicate key entries because mock doesn't have that infrastructure
		// error is on the db layer
		// fmt.Sprintf(`
		// 	# duplicate email
		// 	info: {
		// 		email: "%s",
		// 		password:"password"
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
		assert.Contains(response, "errors", fmt.Sprintf("Case [%v]: %s", i+1, msgNoError))
	}

}
