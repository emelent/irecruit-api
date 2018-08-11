package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestRecruitList(t *testing.T) {
	handler := createLoadedGqlHandler()

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

func TestCreateRecruitValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createRecruit"
	account := moc.Accounts[len(moc.Accounts)-2]
	query := fmt.Sprintf(`
		mutation{
			%s(
				account_id: "%s",
				info: {
					province: KWAZULU_NATAL,
					city: "Durban",
					gender: "male",
					disability: "",
					vid1_url: "http://google.com",
					vid2_url: "http://youtube.com",
					phone: "0123456789",
					email: "dude@gmail.com"
				}
			){
				name,
				surname
			}
		}
	`, method, account.ID.Hex())

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	// make assertions
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	dataPortion := assertGqlData(method, response, assert)
	assert.NotNil(dataPortion[method], msgInvalidResult)
	resultRecruit, rOk := dataPortion[method].(map[string]interface{})
	assert.True(rOk, fmt.Sprintf("Invalid data[\"%s\"] response type.", method))
	if rOk {
		assert.Equal(account.Name, resultRecruit["name"].(string), msgInvalidResult)
	}
}

func TestCreateRecruitInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "createRecruit"
	account := moc.Accounts[len(moc.Accounts)-2]
	queryFormat := `
		mutation{
			%s(
				%s
			){
				name,
				surname
			}
		}
	`

	input := []string{
		fmt.Sprintf(`
			#case 2, invalid RecruitDetails
			account_id: "%s",
			info: {
			}
		`, account.ID.Hex()),
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
		assert.NotContains(response, "data", fmt.Sprintf("Case [%v]: %s", i+1, msgUnexpectedData))
	}
}
