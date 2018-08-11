package models

import (
	"strings"

	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// QA QuestionAnswer db model
type QA struct {
	QuestionID bson.ObjectId `json:"question_id" bson:"question_id"`
	Answer     string        `json:"answer" bson:"answer"`
}

// Recruit db model
type Recruit struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Province   string        `json:"province" bson:"province"`
	City       string        `json:"city" bson:"city"`
	Gender     string        `json:"gender" bson:"gender"`
	Disability string        `json:"disability" bson:"disability"`
	Vid1Url    string        `json:"vid1_url" bson:"vid1_url"`
	Vid2Url    string        `json:"vid2_url" bson:"vid2_url"`
	Phone      string        `json:"phone" bson:"phone"`
	Email      string        `json:"email" bson:"email"`
	// Qa1        QA            `json:"qa1" bson:"qa1"`
	// Qa2        QA            `json:"qa2" bson:"qa2"`
}

//OK validates Recruit fields
func (r *Recruit) OK() error {
	if r.Province == "" {
		return er.NewInvalidFieldError("province")
	}
	if r.City == "" {
		return er.NewInvalidFieldError("city")
	}
	r.Gender = strings.ToLower(r.Gender)
	if !(r.Gender == "male" || r.Gender == "female") {
		return er.NewInvalidFieldError("gender")
	}

	if r.Phone == "" {
		return er.NewInvalidFieldError("phone")
	}
	if r.Email == "" {
		return er.NewInvalidFieldError("email")
	}
	return nil
}
