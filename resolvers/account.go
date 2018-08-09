package resolvers

import (
	"context"
	"log"

	er "../errors"
	mware "../middleware"
	models "../models"
	utils "../utils"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

const accountsCollection = "accounts"

type accountDetails struct {
	Email    string
	Password string
	Name     string
	Surname  string
}

type accountResolver struct {
	a *models.Account
}

func (r *accountResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

func (r *accountResolver) Email() string {
	return r.a.Email
}

func (r *accountResolver) Name() string {
	return r.a.Name
}

func (r *accountResolver) Surname() string {
	return r.a.Surname
}

func (r *accountResolver) AccessLevel() int {
	return r.a.AccessLevel
}

func (r *accountResolver) HunterID() *graphql.ID {
	if r.a.HunterID == nil {
		return nil
	}

	id := graphql.ID(r.a.HunterID.Hex())
	return &id
}

func (r *accountResolver) RecruitID() *graphql.ID {
	if r.a.RecruitID == nil {
		return nil
	}
	id := graphql.ID(r.a.RecruitID.Hex())
	return &id
}

type tokensResolver struct {
	refresh string
	access  string
}

func (r *tokensResolver) AccessToken() string {
	return r.access
}

func (r *tokensResolver) RefreshToken() string {
	return r.refresh
}

// Accounts resolves accounts(name: String) query
func (r *RootResolver) Accounts() ([]*accountResolver, error) {
	defer r.crud.CloseCopy()
	rawAccounts, err := r.crud.FindAll(accountsCollection, nil)
	results := make([]*accountResolver, 0)
	for _, r := range rawAccounts {
		account := transformAccount(r)
		results = append(results, &accountResolver{&account})
	}
	return results, err
}

// CreateAccount resolves the query of the same name
func (r *RootResolver) CreateAccount(ctx context.Context, args struct{ Info *accountDetails }) (*tokensResolver, error) {
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
	genericErr := "Failed to create Account"

	// validate account data
	err := account.OK()
	if err != nil {
		return nil, err
	}

	// hash the new password
	account.HashPassword()

	// store account in db
	err = r.crud.Insert(accountsCollection, account)
	if err != nil {
		log.Println("Failed to create Account =>", err)
		return nil, er.NewInternalError(genericErr)
	}

	// create refresh token
	id := account.ID.Hex()
	refresh, err := utils.CreateRefreshToken(id)
	if err != nil {
		log.Println("Failed to create refresh token =>", err)
		return nil, er.NewGenericError()
	}

	// access token
	ua := ctx.Value(mware.UaKey).(string)
	access, err := utils.CreateAccessToken(id, ua)
	if err != nil {
		log.Println("Failed to create access token =>", err)
		return nil, er.NewGenericError()
	}

	// create token manager
	tokenMgr := models.TokenManager{
		ID:           bson.NewObjectId(),
		AccountID:    account.ID,
		Tokens:       []string{access},
		RefreshToken: refresh,
		MaxTokens:    5,
	}
	err = r.crud.Insert(tokenMgrCollection, tokenMgr)
	if err != nil {
		log.Println("Failed to create TokenManager", err)
		return nil, er.NewInternalError(genericErr)
	}

	return &tokensResolver{refresh: refresh, access: access}, nil
}

// RemoveAccount removes an account
func (r *RootResolver) RemoveAccount(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()
	genericErr := "Failed to remove account."
	idStr := string(args.ID)
	if !bson.IsObjectIdHex(idStr) {
		return nil, er.NewInternalError(genericErr)
	}
	id := bson.ObjectIdHex(idStr)

	// check if there's an account with that id
	_, err := r.crud.FindOne(accountsCollection, &bson.M{"_id": id})
	if err != nil {
		return nil, er.NewInternalError(genericErr)
	}

	// delete the account
	err = r.crud.DeleteID(accountsCollection, id)
	if err != nil {
		log.Println("Failed to delete Account =>", err)
		return nil, er.NewGenericError()
	}

	// find the account's token manager
	rawTokenMgr, err := r.crud.FindOne(tokenMgrCollection, &bson.M{"account_id": id})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	// delete the account's token manager
	tokenMgr := transformTokenManager(rawTokenMgr)
	err = r.crud.DeleteID(tokenMgrCollection, tokenMgr.ID)
	if err != nil {
		log.Println("Failed to delete TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	msg := "Account successfully removed."
	return &msg, nil
}
