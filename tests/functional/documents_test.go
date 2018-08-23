package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateDocumentValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createDocument"
	url := "http://yurp.com"
	ownerID := moc.Recruits[0].ID.Hex()
	query := fmt.Sprintf(`
		mutation{
			%s(doc_type:QUALIFICATION, owner_type: RECRUIT,owner_id: "%s", url: "%s"){
				id
				url
				owner_id
			}
		}
	`, method, ownerID, url)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion, dOk := response["data"].(map[string]interface{})
	resultDocument, rOk := dataPortion[method].(map[string]interface{})

	//make assertions
	assert.NotContains(response, "errors", msgUnexpectedError)
	assert.Contains(response, "data", msgInvalidResponse)
	assert.Contains(response["data"], method, msgMissingResponseData)
	assert.True(dOk, msgInvalidResponseType)
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))

	if rOk {
		assert.Contains(resultDocument, "id", msgMissingResponseData)

		assert.Contains(resultDocument, "url", msgMissingResponseData)
		assert.Equal(resultDocument["url"], url, msgInvalidResult)

		assert.Contains(resultDocument, "owner_id", msgMissingResponseData)
		assert.Equal(resultDocument["owner_id"], ownerID, msgInvalidResult)
	}
}

func TestCreateDocumentInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createDocument"
	queryFormat := `
		mutation{
			%s(%s){
				id
				document
			}
		}
	`

	// invalid inputs
	url := "http://yurp.com"
	ownerID := moc.Recruits[0].ID.Hex()

	input := []string{
		fmt.Sprintf(`
			# case 1 no owner_id
			url: "%s",
			owner_type: RECRUIT,
			doc_type: QUALIFICATION,
		`, url),
		fmt.Sprintf(`
			# case 2 no url
			owner_id: "%s",
			owner_type: RECRUIT,
			doc_type: QUALIFICATION,
		`, ownerID),
		fmt.Sprintf(`
			# case 3 no owner_type
			owner_id: "%s",
			url: "%s",
			doc_type: QUALIFICATION,
		`, ownerID, url),
		fmt.Sprintf(`
			# case 4 no doc_type
			owner_id: "%s",
			url: "%s",
			owner_type: RECRUIT,
		`, ownerID, url),
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
