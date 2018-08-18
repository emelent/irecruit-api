package functionaltests

import (
	"fmt"
	"testing"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
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
					email: "dude@gmail.com",
					qa1_question: "What's up?",
					qa1_answer: "Nothing much.",
					qa2_question: "You good though?",
					qa2_answer: "You know it."
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
			#case 1, no RecruitDetails
			account_id: "%s",
			info: {
			}
		`, account.ID.Hex()),
		`
			#case 2, no account_id
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`,
		`
			#case 2, invalid account_id
			account_id: "sadf",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`,
		fmt.Sprintf(`
			#case 3, missing province
			account_id: "%s",
			info: {
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 4, missing city
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 5, missing gender
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 6, missing disability
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 7, missing phone
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 8, missing email
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 9, missing qa1_question
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 10, missing qa1_answer
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa2_question: "You good though?",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 11, missing qa2_question
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_answer: "You know it."				
			}
		`, account.ID.Hex()),
		fmt.Sprintf(`
			#case 12, missing qa2_answer
			account_id: "%s",
			info: {
				province: KWAZULU_NATAL,
				city: "Durban",
				gender: "male",
				disability: "",
				vid1_url: "http://google.com",
				vid2_url: "http://youtube.com",
				phone: "0123456789",
				email: "dude@gmail.com",
				qa1_question: "What's up?",
				qa1_answer: "Nothing much.",
				qa2_question: "You good though?",			
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
	}
}

func TestRemoveRecruitValid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "removeRecruit"
	query := fmt.Sprintf(`
		mutation{
			%s(id:"%s")
		}
	`, method, moc.Recruits[0].ID.Hex())

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
		assert.Equal(result, "Recruit successfully removed.", msgInvalidResult)
	}
}

func TestRemoveRecruitInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "removeRecruit"
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
