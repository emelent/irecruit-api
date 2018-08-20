package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// Question model
type Question struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	IndustryID bson.ObjectId `json:"industry_id" bson:"industry_id"`
	Question   string        `json:"question" bson:"question"`
}

// OK validate Question model
func (i *Question) OK() error {
	if i.Question == "" {
		return er.InvalidField("question")
	}
	if i.IndustryID == "" {
		return er.InvalidField("industry_id")
	}

	return nil
}
