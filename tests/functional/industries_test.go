package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestIndustryList(t *testing.T) {
	handler := createLoadedGqlHandler()

	//prepare request
	method := "industries"
	query := fmt.Sprintf(`query{%s{id,name}}`, method)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultIndustries, rOk := dataPortion[method].([]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	if rOk {
		assert.Len(resultIndustries, len(moc.Industries), msgInvalidResultCount)
	}
}

func TestCreateIndustryValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createIndustry"
	name := "Test Industry"
	query := fmt.Sprintf(`
		mutation{
			%s(name: "%s"){
				id
				name
			}
		}
	`, method, name)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultIndusty, rOk := dataPortion[method].(map[string]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))

	if rOk {
		assert.Contains(resultIndusty, "id", msgMissingResponseData)
		assert.Contains(resultIndusty, "name", msgMissingResponseData)
		assert.Equal(resultIndusty["name"], name, msgInvalidResult)
	}
}

func TestCreateIndustryInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createIndustry"
	queryFormat := `
		mutation{
			%s(%s){
				id
				name
			}
		}
	`

	// invalid inputs
	input := []string{
		`
			# case 1 no name
		`,
		` 
			# case 2 invalid name
			name: ""
		`,
		fmt.Sprintf(`	# case 3 duplicate name
			name: "%s"
		`, moc.Industries[0].Name),
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
