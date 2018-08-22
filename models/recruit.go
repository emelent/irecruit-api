package models

import (
	"strings"
	"time"

	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformRecruit transforms interface into Recruit model
func TransformRecruit(in interface{}) Recruit {
	var recruit Recruit
	switch v := in.(type) {
	case bson.M:
		recruit.ID = v["_id"].(bson.ObjectId)
		recruit.BirthYear = v["birth_year"].(int32)
		recruit.Province = v["province"].(string)
		recruit.City = v["city"].(string)
		recruit.Gender = v["gender"].(string)
		recruit.Disability = v["disability"].(string)
		recruit.Vid1Url = v["vid1_url"].(string)
		recruit.Vid2Url = v["vid2_url"].(string)
		recruit.Phone = v["phone"].(string)
		recruit.Email = v["email"].(string)
		recruit.Qa1 = TransformQA(v["qa1"])
		recruit.Qa2 = TransformQA(v["qa2"])

	case Recruit:
		recruit = v
	}
	return recruit
}

// -----------------
// Model
// -----------------

// Recruit db model
type Recruit struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	BirthYear  int32         `json:"birth_year" bson:"birth_year"`
	Province   string        `json:"province" bson:"province"`
	City       string        `json:"city" bson:"city"`
	Gender     string        `json:"gender" bson:"gender"`
	Disability string        `json:"disability" bson:"disability"`
	Vid1Url    string        `json:"vid1_url" bson:"vid1_url"`
	Vid2Url    string        `json:"vid2_url" bson:"vid2_url"`
	Phone      string        `json:"phone" bson:"phone"`
	Email      string        `json:"email" bson:"email"`
	Qa1        QA            `json:"qa1" bson:"qa1"`
	Qa2        QA            `json:"qa2" bson:"qa2"`
}

//OK validates Recruit fields
func (r *Recruit) OK() error {
	if r.Province == "" {
		return er.InvalidField("province")
	}
	if r.City == "" {
		return er.InvalidField("city")
	}
	r.Gender = strings.ToLower(r.Gender)
	if !(r.Gender == "male" || r.Gender == "female") {
		return er.InvalidField("gender")
	}

	if r.Phone == "" {
		return er.InvalidField("phone")
	}
	if r.Email == "" {
		return er.InvalidField("email")
	}

	if r.Qa1.Question == "" {
		return er.InvalidField("qa1.question")
	}

	if r.Qa1.Answer == "" {
		return er.InvalidField("qa1.answer")
	}

	if r.Qa2.Question == "" {
		return er.InvalidField("qa2.question")
	}

	if r.Qa2.Answer == "" {
		return er.InvalidField("qa2.answer")
	}

	if r.BirthYear < 1900 || r.BirthYear >= int32(time.Now().Year()) {
		return er.InvalidField("birth_year")
	}

	r.Gender = strings.ToUpper(r.Gender)
	return nil
}
