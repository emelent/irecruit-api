package functionaltests

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

func TestIndustryList(t *testing.T) {

	//prepare request
	handler := createLoadedGqlHandler()
	method := "industries"
	query := fmt.Sprintf(`query{%s{name}}`, method)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	assert := assert.New(t)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	/* expected result
	list := make([]interface{}, 0)
	for _, i := range moc.Industries {
		list = append(list, map[string]interface{}{
			"name": strings.ToLower(i.Name),
		})
	}

	expect := map[string]interface{}{
		"data": map[string]interface{}{
			method: list,
		},
	}
	assert.Equal(expect, response, msgInvalidResult)
}

func TestCreateIndustryValid(t *testing.T) {
	handler := createLoadedGqlHandler()
	assert := assert.New(t)

	// prepare request
	method := "createIndustry"
	name := "Test Industry"
	query := fmt.Sprintf(`
		mutation{
			%s(name: "%s"){
				name
			}
		}
	`, method, name)

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	// expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			method: map[string]interface{}{
				"name": strings.ToLower(name),
			},
		},
	}
	assert.Equal(expected, response, msgInvalidResult)
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

func TestRemoveIndustryValid(t *testing.T) {
	handler := createLoadedGqlHandler()
	assert := assert.New(t)

	// prepare request
	method := "removeIndustry"
	query := fmt.Sprintf(`
		mutation{
			%s(id:"%s")
		}
	`, method, moc.Industries[0].ID.Hex())

	// request and respond
	response, err := gqlRequestAndRespond(handler, query, nil)

	//process response
	if err != nil {
		assert.Fail("Failed to process response:", err)
	}

	// expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			method: "Industry successfully removed.",
		},
	}
	assert.Equal(expected, response, msgInvalidResult)
}

func TestRemoveIndustryInvalid(t *testing.T) {
	handler := createLoadedGqlHandler()

	// prepare request
	method := "removeIndustry"
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
