package functionaltests

import (
	"fmt"
	"testing"
	"time"

	moc "../../mocks"
	"github.com/stretchr/testify/assert"
)

// tests edit on AccountEditor
func TestEditAccountEditor(t *testing.T) {
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

// tests edit on RecruitEditor
func TestEditRecruitEditor(t *testing.T) {
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

func TestRecruitEditorUpdateQAs(t *testing.T) {
	assert := assert.New(t)
	crud := moc.NewLoadedCRUD()
	handler := createGqlHandler(crud)

	// login as recruit user
	token, _ := login(crud, getRecruitUserAccount().ID, "none")

	// prep data
	qa1Question := moc.Questions[0]
	qa2Question := moc.Questions[0]
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
