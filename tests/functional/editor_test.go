package functionaltests

import (
	"fmt"
	"testing"
	"time"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

// tests that AccountEditor.UpdateAccount updates the current  account
func TestAccountEditor_UpdateAccount(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as plain user
	token, _ := login(crud, getPlainUserAccount().ID, "none")

	// prepare query
	name := "Newa"
	surname := "Fields"
	email := "newa@gmail.com"
	query := fmt.Sprintf(`
		mutation {
			edit(token: "%s"){
				... on AccountEditor{
					updateAccount(info: {
						name: "%s",
						surname: "%s",
						email: "%s"
					}){
						name
						surname
						email
					}
				}
			}
		}
	`, token, name, surname, email)

	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"updateAccount": map[string]interface{}{
					"name":    name,
					"surname": surname,
					"email":   email,
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

func TestAccountEditor_RemoveAccount(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as plain user
	token, _ := login(crud, getPlainUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on AccountEditor{
					removeAccount
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeAccount": "Account successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

func TestAccountEditor_CreateRecruit(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as plain user
	token, _ := login(crud, getPlainUserAccount().ID, "none")

	// prepare data
	phone := "082 345 6789"
	email := "newa@gmail.com"
	province := "LIMPOPO"
	gender := "FEMALE"
	city := "Polokwane"
	disability := "Blind"
	vid1URL := "http://google.com"
	vid2URL := "http://youtube.com"
	qa1Question := moc.Questions[0]
	qa2Question := moc.Questions[1]
	qa1Answer := "Interesting... very interesting."
	qa2Answer := "No idea."
	birthYear := 1987

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on AccountEditor{
					createRecruit(info:{
						phone: "%s",
						email: "%s",
						province: %s,
						city: "%s",
						gender: %s,
						disability: "%s",
						vid1_url: "%s",
						vid2_url: "%s",
						birth_year: %v,
						qa1_question_id: "%s",
						qa1_answer: "%s",
						qa2_question_id: "%s",
						qa2_answer: "%s",
					}){
						phone
						email
						province
						city
						gender
						disability
						vid1_url
						vid2_url
						age
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
	`, token, phone, email, province, city, gender, disability,
		vid1URL, vid2URL, birthYear, qa1Question.ID.Hex(), qa1Answer,
		qa2Question.ID.Hex(), qa2Answer,
	)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"createRecruit": map[string]interface{}{
					"phone":      phone,
					"email":      email,
					"province":   province,
					"city":       city,
					"gender":     gender,
					"disability": disability,
					"vid1_url":   vid1URL,
					"vid2_url":   vid2URL,
					"age":        float64(time.Now().Year() - birthYear),
					"qa1": map[string]interface{}{
						"question": qa1Question.Question,
						"answer":   qa1Answer,
					},
					"qa2": map[string]interface{}{
						"question": qa2Question.Question,
						"answer":   qa2Answer,
					},
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests that RecruitEditor.UpdateRecruit updates Recruit
func TestRecruitEditor_UpdateRecruit(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as recruit user
	token, _ := login(crud, getRecruitUserAccount().ID, "none")

	// prep data
	phone := "082 345 6789"
	email := "newa@gmail.com"
	province := "LIMPOPO"
	gender := "FEMALE"
	city := "Polokwane"
	disability := "Blind"
	vid1URL := "http://google.com"
	vid2URL := "http://youtube.com"
	birthYear := 1987

	// prepare query
	query := fmt.Sprintf(`
		mutation {
			edit(token: "%s"){
				... on RecruitEditor{
					updateRecruit(info: {
						phone: "%s",
						email: "%s",
						province: %s,
						city: "%s",
						gender: %s,
						disability: "%s",
						vid1_url: "%s",
						vid2_url: "%s",
						birth_year: %v,
					}){
						phone
						email
						province
						city
						gender
						disability
						vid1_url
						vid2_url
						age
					}
				}
			}
		}
	`, token, phone, email, province, city, gender, disability,
		vid1URL, vid2URL, birthYear,
	)

	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"updateRecruit": map[string]interface{}{
					"phone":      phone,
					"email":      email,
					"province":   province,
					"city":       city,
					"gender":     gender,
					"disability": disability,
					"vid1_url":   vid1URL,
					"vid2_url":   vid2URL,
					"age":        float64(time.Now().Year() - birthYear),
				},
			},
		},
	}
	assert.Equal(expected, response, msgInvalidResult)
}

// tests that RecruitEditor.UpdateQAs updates the current Recruit profiles QAs
func TestRecruitEditor_UpdateQAs(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as recruit user
	token, _ := login(crud, getRecruitUserAccount().ID, "none")

	// prep data
	qa1Question := moc.Questions[0]
	qa2Question := moc.Questions[1]
	qa1Answer := "Interesting... very interesting."
	qa2Answer := "No idea."

	// prep query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s", enforce: RECRUIT){
				... on RecruitEditor{
					updateQAs(
						qa1:{
							question_id: "%s",
							answer: "%s"
						}, 
						qa2:{
							question_id: "%s",
							answer: "%s"
						}
					){
						question
						answer
					}
				}
			}
		}
	`, token, qa1Question.ID.Hex(), qa1Answer, qa2Question.ID.Hex(), qa2Answer)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"updateQAs": []interface{}{
					map[string]interface{}{
						"question": qa1Question.Question,
						"answer":   qa1Answer,
					},
					map[string]interface{}{
						"question": qa2Question.Question,
						"answer":   qa2Answer,
					},
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

func TestRecruitEditor_RemoveRecruit(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as recruit user
	token, _ := login(crud, getRecruitUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on RecruitEditor{
					removeRecruit
				}
			}
		}
	`, token)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeRecruit": "Recruit successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResult)
}

// tests that edit mutation cannot be accessed without a valid token
func TestEditWithInvalidToken(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on Editor{}
			}
		}
	`, "bad_token")

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)
	assert.Contains(response, "errors", msgNoError)
}

func TestSysEditor_RemoveAccount(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					removeAccount(id: "%s")
				}
			}
		}
	`, token, getPlainUserAccount().ID.Hex())

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeAccount": "Account successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_RemoveRecruit(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					removeRecruit(id: "%s")
				}
			}
		}
	`, token, getRecruitUserAccount().RecruitID.Hex())

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeRecruit": "Recruit successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}
func TestSysEditor_RemoveIndustry(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					removeIndustry(id: "%s")
				}
			}
		}
	`, token, moc.Industries[0].ID.Hex())

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeIndustry": "Industry successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_RemoveQuestion(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					removeQuestion(id: "%s")
				}
			}
		}
	`, token, moc.Questions[0].ID.Hex())

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeQuestion": "Question successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_RemoveDocument(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					removeDocument(id: "%s")
				}
			}
		}
	`, token, moc.Documents[0].ID.Hex())

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"removeDocument": "Document successfully removed.",
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_CreateQuestion(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prep data
	industryID := moc.Industries[0].ID.Hex()
	question := "Have you ever tasted stuffed socks?"

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					createQuestion(industry_id: "%s", question: "%s"){
						question
						industry_id
					}
				}
			}
		}
	`, token, industryID, question)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"createQuestion": map[string]interface{}{
					"question":    question,
					"industry_id": industryID,
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_CreateIndustry(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	name := "megastructures"

	// prepare query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					createIndustry(name: "%s"){
						name
					}
				}
			}
		}
	`, token, name)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prep expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"createIndustry": map[string]interface{}{
					"name": name,
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}

func TestSysEditor_UpdateIndustry(t *testing.T) {
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)
	assert := assert.New(t)

	// login as sys
	token, _ := login(crud, getSysUserAccount().ID, "none")

	// prep data
	name := "schwartz"
	industryID := moc.Industries[0].ID.Hex()

	// prep query
	query := fmt.Sprintf(`
		mutation{
			edit(token: "%s"){
				... on SysEditor{
					updateIndustry(id: "%s", name: "%s"){
						id
						name
					}
				}
			}
		}
	`, token, industryID, name)

	// request
	response, err := gqlRequestAndRespond(handler, query, nil)
	failOnError(assert, err)

	// prepare expected
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"edit": map[string]interface{}{
				"updateIndustry": map[string]interface{}{
					"id":   industryID,
					"name": name,
				},
			},
		},
	}

	assert.Equal(expected, response, msgInvalidResponse)
}
