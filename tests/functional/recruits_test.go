package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestRecruitList(t *testing.T) {
	//prepare handler
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	//prepare request
	method := "recruits"
	query := fmt.Sprintf(`query{%s{id,name,surname,gender,phone,email,province,city}}`, method)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultRecruits, rOk := dataPortion[method].([]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	assert.Len(resultRecruits, len(moc.Recruits), msgInvalidResultCount)
}
