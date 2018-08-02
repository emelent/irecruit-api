package models

import (
	e "../errors"
	"gopkg.in/mgo.v2/bson"
)

// Token model
type Token struct {
	Signature string `json:"signature" bson:"signature"`
	Device    string `json:"device" bson:"device"`
	IP        string `json:"ip" bson:"ip"`
	UserAgent string `json:"user_agent" bson:"user_agent"`
}

//OK validates Token fields
func (t *Token) OK() error {
	if t.Signature == "" {
		return e.NewMissingFieldError("Signature")
	}
	return nil
}

// TokenManager model
type TokenManager struct {
	AccountID    bson.ObjectId `json:"account_id" bson:"account_id"`
	Tokens       []Token       `json:"tokens" bson:"tokens"`
	RefreshToken Token         `json:"refresh_token" bson:"refresh_token"`

	maxTokens int
}

// OK validates token manager
func (tm *TokenManager) OK() error {
	return nil
}
