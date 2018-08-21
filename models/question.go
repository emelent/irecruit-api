package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformQuestion transforms interface into Question model
func TransformQuestion(in interface{}) Question {
	var question Question
	switch v := in.(type) {
	case bson.M:
		question.ID = v["_id"].(bson.ObjectId)
		question.IndustryID = v["industry_id"].(bson.ObjectId)
		question.Question = v["question"].(string)

	case Question:
		question = v
	}

	return question
}

// -----------------
// Model
// -----------------

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
