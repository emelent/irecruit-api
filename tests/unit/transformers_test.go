package unittests

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	models "../../models"
	"github.com/stretchr/testify/assert"
)

func TestAccountTransformer(t *testing.T) {
	assert := assert.New(t)

	b := bson.M{
		"_id":          bson.NewObjectId(),
		"email":        "hey@gmail.com",
		"name":         "Jamal",
		"surname":      "Ccom",
		"password":     "hey@gmail.com",
		"access_level": 0,
		"hunter_id":    bson.NewObjectId(),
		"recruit_id":   bson.NewObjectId(),
	}

	expected := models.Account{
		ID:          b["_id"].(bson.ObjectId),
		Email:       b["email"].(string),
		Name:        b["name"].(string),
		Surname:     b["surname"].(string),
		Password:    b["password"].(string),
		AccessLevel: b["access_level"].(int),
		HunterID:    b["hunter_id"].(bson.ObjectId),
		RecruitID:   b["recruit_id"].(bson.ObjectId),
	}

	assert.Equal(expected, models.TransformAccount(b))
}

func TestDocumentTransformer(t *testing.T) {
	assert := assert.New(t)

	b := bson.M{
		"_id":        bson.NewObjectId(),
		"url":        "http://google.com",
		"doc_type":   "QUALIFICATION",
		"owner_type": "RECRUIT",
		"owner_id":   bson.NewObjectId(),
	}

	expected := models.Document{
		ID:        b["_id"].(bson.ObjectId),
		URL:       b["url"].(string),
		OwnerType: b["owner_type"].(string),
		OwnerID:   b["owner_id"].(bson.ObjectId),
		DocType:   b["doc_type"].(string),
	}

	assert.Equal(expected, models.TransformDocument(b))
}
