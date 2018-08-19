package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// Document model
type Document struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	URL       string        `json:"url" bson:"url"`
	DocType   string        `json:"doc_type" bson:"doc_type"`
	OwnerType string        `json:"owner_type"  bson:"owner_type"`
	OwnerID   bson.ObjectId `json:"owner_id" bson:"owner_id"`
}

// OK validates fields of document model
func (d *Document) OK() error {
	if d.URL == "" {
		return er.NewInvalidFieldError("url")
	}
	if d.DocType == "" {
		return er.NewInvalidFieldError("doc_type")
	}
	if d.OwnerType == "" {
		return er.NewInvalidFieldError("owner_type")
	}
	return nil
}
