package functionaltests

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestQuestionList(t *testing.T) {
	handler := createLoadedGqlHandler()

	//prepare request
	method := "questions"
	query := fmt.Sprintf(`query{%s{id,question,industry_id}}`, method)

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
		assert.Len(resultQuestions, len(moc.Questions), msgInvalidResultCount)
	}
}

func TestCreateQuestionValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createQuestion"
	question := "What's that?"
	query := fmt.Sprintf(`
		mutation{
			%s(industry_id: "%s", question: "%s"){
				id
				question
			}
		}
	`, method, moc.Industries[0].ID.Hex(), question)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultQuestion, rOk := dataPortion[method].(map[string]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))

	if rOk {
		assert.Contains(resultQuestion, "id", msgMissingResponseData)
		assert.Contains(resultQuestion, "question", msgMissingResponseData)
		assert.Equal(resultQuestion["question"], question, msgInvalidResult)
	}
}

func TestCreateQuestionInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createQuestion"
	queryFormat := `
		mutation{
			%s(%s){
				id
				question
			}
		}
	`

	// invalid inputs
	input := []string{
		fmt.Sprintf(`
			# case 1 no question
			industry_id: "%s"
		`, moc.Industries[0].ID.Hex()),
		` 
			# case 2 no industry_id
			question: "What's up?"
		`,
		`
			# case 3 invalid industry_id
			industry_id: "43",
			question: "Dude?"
		`,
		fmt.Sprintf(`
			# case 4 invalid question
			industry_id: "%s",
			question: ""
		`, moc.Industries[0].ID.Hex()),
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

func TestRemoveQuestionValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "removeQuestion"
	query := fmt.Sprintf(`
		mutation{
			%s(id:"%s")
		}
	`, method, moc.Questions[0].ID.Hex())

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

	if rOk {
		assert.Equal(result, "Question successfully removed.", msgInvalidResult)
	}
}

func TestRemoveQuestionInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "removeQuestion"
	queryFormat := `
		mutation{
			%s(%s)
		}
	`

	// invalid inputs
	input := []string{
		`
			# case 1 no id
		`,
		`
			# case 2 invalid id
			id: "id"
		`,
		fmt.Sprintf(`
			# case 3 non-existent id
			id: "%s"
		`, bson.NewObjectId()),
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
