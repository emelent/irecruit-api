package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformQA transforms interface into QA model
func TransformQA(in interface{}) QA {
	var qa QA
	switch v := in.(type) {
	case bson.M:
		qa.Question = v["question"].(string)
		qa.Answer = v["answer"].(string)
	case QA:
		qa = v
	}
	return qa
}

// -----------------
// Model
// -----------------

// QA QuestionAnswer db model
type QA struct {
	Question string `json:"question" bson:"question"`
	Answer   string `json:"answer" bson:"answer"`
}

// OK validates QA fields
func (q *QA) OK() error {
	if q.Question == "" {
		return er.InvalidField("question")
	}

	if q.Answer == "" {
		return er.InvalidField("answer")
	}
	return nil
}
