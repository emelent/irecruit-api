package functionaltests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

// tests view on SysViewer
func TestViewSysViewer(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as sys admin
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on SysViewer{
					accounts{
						name
					}
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// expected
	list := make([]map[string]interface{}, 0)
	for _, acc := range moc.Accounts {
		list = append(list, map[string]interface{}{
			"name": acc.Name,
		})
	}
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"accounts": list,
			},
		},
	}

	// use strings because it's easier to construct
	// the expected result than with actual types
	expectedStr := fmt.Sprintf("%s", expected)
	actualStr := fmt.Sprintf("%s", response)

	assert.Equal(expectedStr, actualStr, msgInvalidResult)
}

// tests that View returns an error on an invalid token
func TestViewWithInvalidToken(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on Viewer{}
			}
		}
	`, "")

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)
	assert.Contains(response, "errors", msgNoError)
}

// tests view on RecruitViewer
func TestViewRecruitViewer(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as recruit account
	token, _ := login(crud, getRecruitUserAccount().ID, "none")
	recruit := moc.Recruits[0]

	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on RecruitViewer{
					profile{
						id
					}
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"profile": map[string]interface{}{
					"id": recruit.ID.Hex(),
				},
			},
		},
	}

	// use strings because it's easier to construct
	// the expected result than with actual types
	expectedStr := fmt.Sprintf("%s", expected)
	actualStr := fmt.Sprintf("%s", response)

	assert.Equal(expectedStr, actualStr, msgInvalidResult)
}

// tests view on AccountViewer
func TestViewAccountViewer(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login with plain account
	token, _ := login(crud, getPlainUserAccount().ID, "none")
	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on AccountViewer{
					is_hunter
					is_recruit
					checkPassword(password: "password")
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"is_hunter":     false,
				"is_recruit":    false,
				"checkPassword": true,
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests view on Viewer
func TestViewViewer(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login with plain account
	account := getPlainUserAccount()
	token, _ := login(crud, account.ID, "none")
	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on Viewer{
					id
					name
					surname
					email
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"id":      account.ID.Hex(),
				"name":    account.Name,
				"surname": account.Surname,
				"email":   strings.ToLower(account.Email),
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests valid viewer enforcing
func TestViewEnforceValid(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// get access tokens
	sysToken, _ := login(crud, getSysUserAccount().ID, "none")
	recruitToken, _ := login(crud, getRecruitUserAccount().ID, "none")

	// prepare query
	queryFormat := `query{
		view(%s){
			... on Viewer{}
		}
	}`

	// valid  enforce cases
	input := []string{
		fmt.Sprintf(`
			# case 1 enforce SYSTEM
			token: "%s", enforce: SYSTEM
		`, sysToken),
		fmt.Sprintf(`
			# case 2 enforce RECRUIT
			token: "%s", enforce: RECRUIT
		`, recruitToken),
		fmt.Sprintf(`
			# case 3 system account enforce ACCOUNT
			token: "%s", enforce: ACCOUNT
		`, sysToken),
		fmt.Sprintf(`
			# case 3 recruit account enforce ACCOUNT
			token: "%s", enforce: ACCOUNT
		`, recruitToken),
	}

	for i, in := range input {
		query := fmt.Sprintf(queryFormat, in)

		// request
		response, err := gqlRequestAndRespond(handler, query, nil)
		failOnError(assert, err)
		assert.NotContains(response, "errors", fmt.Sprintf("Case [%v]: %s", i+1, msgUnexpectedError))
	}
}

// tests invalid viewer enforcing
func TestViewEnforceInvalid(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// get access tokens
	recruitToken, _ := login(crud, getRecruitUserAccount().ID, "none")
	plainToken, _ := login(crud, getPlainUserAccount().ID, "none")

	badToken := "bad to the bone"

	// prepare query
	queryFormat := `query{
		view(%s){
			... on Viewer{}
		}
	}`

	// valid  enforce cases
	input := []string{
		fmt.Sprintf(`
			# case 1 enforce SYSTEM on Non-SysAccount
			token: "%s", enforce: SYSTEM
		`, recruitToken),
		fmt.Sprintf(`
			# case 2 enforce RECRUIT on Non-RecruitAccount
			token: "%s", enforce: RECRUIT
		`, plainToken),
		fmt.Sprintf(`
			# case 3 enforce ACCOUNT on bad token
			token: "%s", enforce: ACCOUNT
		`, badToken),
	}

	for i, in := range input {
		query := fmt.Sprintf(queryFormat, in)

		// request
		response, err := gqlRequestAndRespond(handler, query, nil)
		failOnError(assert, err)
		assert.Contains(response, "errors", fmt.Sprintf("Case [%v]: %s", i+1, msgNoError))
	}
}

// tests that SysViewer.Recruits lists all recruits
func TestViewRecruitList(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys user
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		query {
			view(token: "%s"){
				... on SysViewer{
					recruits{
						id,
						email,
						phone,
						email
						province,
						city,
						disability,
						gender,
						vid1_url,
						vid2_url,
						age,
						qa1{
							question
							answer
						},
						qa2{
							question
							answer
						}

					}
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	year := int32(time.Now().Year())
	recruits := make([]interface{}, 0)
	for _, r := range moc.Recruits {
		recruits = append(recruits, map[string]interface{}{
			"id":         r.ID.Hex(),
			"email":      r.Email,
			"phone":      r.Phone,
			"province":   r.Province,
			"city":       r.City,
			"gender":     r.Gender,
			"disability": r.Disability,
			"vid1_url":   r.Vid1Url,
			"vid2_url":   r.Vid2Url,
			"age":        float64(year - r.BirthYear),
			"qa1": map[string]interface{}{
				"question": r.Qa1.Question,
				"answer":   r.Qa1.Answer,
			},
			"qa2": map[string]interface{}{
				"question": r.Qa2.Question,
				"answer":   r.Qa2.Answer,
			},
		})
	}
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"recruits": recruits,
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests that SysViewer.Questions lists all questions
func TestViewQuestionList(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys user
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on SysViewer{
					questions{
						id
						industry_id
						question
					}
				}
			}
		}
	`, token)

	// request it
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	questions := make([]interface{}, 0)
	for _, question := range moc.Questions {
		questions = append(questions, map[string]interface{}{
			"id":          question.ID.Hex(),
			"industry_id": question.IndustryID.Hex(),
			"question":    question.Question,
		})
	}
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"questions": questions,
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests that SysViewer.Documents lists all documents
func TestViewDocumentsLists(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		query{
			view(token: "%s"){
				... on SysViewer{
					documents{
						id
						url
						doc_type
						owner_type
						owner_id
					}
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	documents := make([]interface{}, 0)
	for _, doc := range moc.Documents {
		documents = append(documents, map[string]interface{}{
			"id":         doc.ID.Hex(),
			"owner_id":   doc.OwnerID.Hex(),
			"owner_type": doc.OwnerType,
			"doc_type":   doc.DocType,
			"url":        doc.URL,
		})
	}
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"view": map[string]interface{}{
				"documents": documents,
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}
