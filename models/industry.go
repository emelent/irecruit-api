package models

import (
	"strings"

	er "../errors"
	"gopkg.in/mgo.v2/bson"
)

// Industry model
type Industry struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

// OK validate Industry model
func (i *Industry) OK() error {
	if i.Name == "" {
		return er.NewInvalidFieldError("name")
	}

	i.Name = strings.ToLower(i.Name)
	return nil
}
