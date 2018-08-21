package models

import (
	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformDocument transforms interface into Document model
func TransformDocument(in interface{}) Document {
	var document Document
	switch v := in.(type) {
	case bson.M:
		document.ID = v["_id"].(bson.ObjectId)
		document.OwnerID = v["owner_id"].(bson.ObjectId)
		document.URL = v["url"].(string)
		document.DocType = v["doc_type"].(string)
		document.OwnerType = v["owner_type"].(string)

	case Document:
		document = v
	}

	return document
}

// -----------------
// Model
// -----------------

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
		return er.InvalidField("url")
	}
	if d.DocType == "" {
		return er.InvalidField("doc_type")
	}
	if d.OwnerType == "" {
		return er.InvalidField("owner_type")
	}
	return nil
}
