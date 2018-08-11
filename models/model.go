package models

import (
	"gopkg.in/mgo.v2/bson"
)

// NullObjectID Serves as place holder for null objectID
var NullObjectID = bson.NewObjectId()

// Model interface
type Model interface {
	OK() error
}
