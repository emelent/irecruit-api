package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// Question model
type Question struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	IndustryID bson.ObjectId `json:"industry_id" bson:"industry_id"`
	Question   string        `json:"name" bson:"name"`
}

// OK validate Question model
func (i *Question) OK() error {
	if i.Question == "" {
		return er.NewInvalidFieldError("question")
	}
	if i.IndustryID == "" {
		return er.NewInvalidFieldError("industry_id")
	}

	return nil
}
