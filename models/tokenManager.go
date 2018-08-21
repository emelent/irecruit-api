package models

import (
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformTokenManager transforms interface into TokenManager model
func TransformTokenManager(in interface{}) TokenManager {
	var tokenMgr TokenManager
	switch v := in.(type) {
	case bson.M:
		tokenMgr.ID = v["_id"].(bson.ObjectId)
		tokenMgr.RefreshToken = v["refresh_token"].(string)
		tokenMgr.MaxTokens = v["max_tokens"].(int)
		tokenMgr.AccountID = v["account_id"].(bson.ObjectId)
		tokenMgr.Tokens = v["tokens"].([]string)

	case TokenManager:
		tokenMgr = v
	}

	return tokenMgr
}

// -----------------
// Model
// -----------------

// TokenManager model
type TokenManager struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	AccountID    bson.ObjectId `json:"account_id" bson:"account_id"`
	Tokens       []string      `json:"tokens" bson:"tokens"`
	RefreshToken string        `json:"refresh_token" bson:"refresh_token"`
	MaxTokens    int           `json:"max_tokens" bson:"max_tokens"`
}

// OK validates token manager
func (tm *TokenManager) OK() error {
	return nil
}
