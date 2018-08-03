package models

import (
	"gopkg.in/mgo.v2/bson"
)

// TokenManager model
type TokenManager struct {
	AccountID    bson.ObjectId `json:"account_id" bson:"account_id"`
	Tokens       []string      `json:"tokens" bson:"tokens"`
	RefreshToken string        `json:"refresh_token" bson:"refresh_token"`
	MaxTokens    int           `json:"max_tokens" bson:"max_tokens"`
}

// OK validates token manager
func (tm *TokenManager) OK() error {
	return nil
}
