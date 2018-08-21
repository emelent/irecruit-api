package models

import (
	"strings"

	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformIndustry transforms interface into Industry model
func TransformIndustry(in interface{}) Industry {
	var industry Industry
	switch v := in.(type) {
	case bson.M:
		industry.ID = v["_id"].(bson.ObjectId)
		industry.Name = v["name"].(string)

	case Industry:
		industry = v
	}

	return industry
}

// -----------------
// Model
// -----------------

// Industry model
type Industry struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

// OK validate Industry model
func (i *Industry) OK() error {
	if i.Name == "" {
		return er.InvalidField("name")
	}

	i.Name = strings.ToLower(i.Name)
	return nil
}
