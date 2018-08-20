package resolvers

import (
	"context"
	"log"

	config "../config"
	er "../errors"
	mware "../middleware"
	models "../models"
	utils "../utils"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Root Resolver methods
// -----------------

// Accounts resolves accounts(name: String) query
func (r *RootResolver) Accounts() ([]*AccountResolver, error) {
	defer r.crud.CloseCopy()
	// get accounts
	rawAccounts, err := r.crud.FindAll(config.AccountsCollection, nil)
	if err != nil {
		log.Println(err)
		return nil, er.Generic()
	}

	// process results
	results := make([]*AccountResolver, 0)
	for _, raw := range rawAccounts {
		account := TransformAccount(raw)
		results = append(results, &AccountResolver{&account})
	}
	return results, err
}

// CreateAccount resolves the query of the same name
func (r *RootResolver) CreateAccount(ctx context.Context, args struct{ Info *accountDetails }) (*TokensResolver, error) {
	defer r.crud.CloseCopy()

	// create account
	account := models.Account{}
	info := args.Info
	account.Name = info.Name
	account.Email = info.Email
	account.Surname = info.Surname
	account.AccessLevel = 0
	account.Password = info.Password
	account.ID = bson.NewObjectId()
	account.HunterID = models.NullObjectID
	account.RecruitID = models.NullObjectID
	genericErr := "Failed to create Account"

	// validate account data
	err := account.OK()
	if err != nil {
		return nil, err
	}

	// hash the new password
	account.HashPassword()

	// store account in db
	err = r.crud.Insert(config.AccountsCollection, account)
	if err != nil {
		log.Println("Failed to create Account =>", err)
		return nil, er.Internal(genericErr)
	}

	// create refresh token
	id := account.ID.Hex()
	refresh, err := utils.CreateRefreshToken(id)
	if err != nil {
		log.Println("Failed to create refresh token =>", err)
		return nil, er.Generic()
	}

	// access token
	ua := ctx.Value(mware.UaKey).(string)
	access, err := utils.CreateAccessToken(id, ua)
	if err != nil {
		log.Println("Failed to create access token =>", err)
		return nil, er.Generic()
	}

	// create token manager
	tokenMgr := models.TokenManager{
		ID:           bson.NewObjectId(),
		AccountID:    account.ID,
		Tokens:       []string{access},
		RefreshToken: refresh,
		MaxTokens:    5,
	}

	// store TokenManager in db
	err = r.crud.Insert(config.TokenManagersCollection, tokenMgr)
	if err != nil {
		log.Println("Failed to create TokenManager", err)
		return nil, er.Internal(genericErr)
	}

	return &TokensResolver{refresh: refresh, access: access}, nil
}

// -----------------
// accountDetails struct
// -----------------
type accountDetails struct {
	Email    string
	Password string
	Name     string
	Surname  string
}

// -----------------
// AccountResolver struct
// -----------------

// AccountResolver resolves account
type AccountResolver struct {
	a *models.Account
}

// ID resolves Account.ID
func (r *AccountResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

// Email resolves Account.Email
func (r *AccountResolver) Email() string {
	return r.a.Email
}

// Name resolves Account.Name
func (r *AccountResolver) Name() string {
	return r.a.Name
}

// Surname resolves Account.Surname
func (r *AccountResolver) Surname() string {
	return r.a.Surname
}

// AccessLevel resolves Account.AccessLevel
func (r *AccountResolver) AccessLevel() int {
	return r.a.AccessLevel
}

// HunterID resolves Account.HunterID
func (r *AccountResolver) HunterID() graphql.ID {
	if r.a.HunterID == models.NullObjectID {
		return graphql.ID("")
	}
	return graphql.ID(r.a.HunterID.Hex())
}

// RecruitID resolves Account.RecruitID
func (r *AccountResolver) RecruitID() graphql.ID {
	if r.a.RecruitID == models.NullObjectID {
		return graphql.ID("")
	}

	return graphql.ID(r.a.RecruitID.Hex())
}

// -----------------
// TokensResolver struct
// -----------------

// TokensResolver resolves Tokens
type TokensResolver struct {
	refresh string
	access  string
}

// AccessToken resolves Token.AccessToken
func (r *TokensResolver) AccessToken() string {
	return r.access
}

// RefreshToken resolves Token.AccessToken
func (r *TokensResolver) RefreshToken() string {
	return r.refresh
}
