package resolvers

import (
	"context"
	"fmt"

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

type failResolver struct {
	message string
}

func (r *failResolver) Error() string {
	return r.message
}

type accountOrFailResolver struct {
	result interface{}
}

func (r *accountOrFailResolver) ToFail() (*failResolver, bool) {
	res, ok := r.result.(*failResolver)
	return res, ok
}

func (r *accountOrFailResolver) ToAccount() (*accountResolver, bool) {
	res, ok := r.result.(*accountResolver)
	return res, ok
}

type tokensOrFailResolver struct {
	result interface{}
}

func (r *tokensOrFailResolver) ToFail() (*failResolver, bool) {
	res, ok := r.result.(*failResolver)
	return res, ok
}

func (r *tokensOrFailResolver) ToTokens() (*tokensResolver, bool) {
	res, ok := r.result.(*tokensResolver)
	return res, ok
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
func (r *RootResolver) Accounts(args struct{ Name *string }) ([]*accountResolver, error) {
	defer r.crud.CloseCopy()
	if args.Name == nil {
		// TODO return all names
	}
	rawAccounts, err := r.crud.FindAll(accountsCollection, nil)
	results := make([]*accountResolver, 0)
	for _, r := range rawAccounts {
		account := transformAccount(r)
		results = append(results, &accountResolver{&account})
	}
	return results, err
}

// CreateAccount resolves the query of the same name
func (r *RootResolver) CreateAccount(ctx context.Context, args struct{ Info *accountDetails }) *tokensOrFailResolver {
	defer r.crud.CloseCopy()

	// create account
	account := models.Account{}
	info := args.Info
	account.Name = info.Name
	account.Email = info.Email
	account.Surname = info.Surname
	account.AccessLevel = 0
	account.SetPassword(info.Password)
	account.ID = bson.NewObjectId()
	genericErr := "Failed to create Account"

	// validate account data
	err := account.OK()
	if err != nil {
		fmt.Println("Account validation failed =>", err)
		return &tokensOrFailResolver{&failResolver{err.Error()}}
	}

	// store account in db
	err = r.crud.Insert(accountsCollection, account)
	if err != nil {
		fmt.Println("Failed to create Account =>", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
	}

	// create refresh token
	id := account.ID.Hex()
	refresh, err := utils.CreateRefreshToken(id)
	if err != nil {
		fmt.Println("Failed to create refresh token =>", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
	}

	// access token
	ua := ctx.Value(mware.UaKey).(string)
	access, err := utils.CreateAccessToken(id, ua)
	if err != nil {
		fmt.Println("Failed to create access token =>", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
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
		fmt.Println("Failed to create TokenManager", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
	}

	result := &tokensResolver{refresh: refresh, access: access}
	return &tokensOrFailResolver{result}
}

// RemoveAccount removes an account
func (r *RootResolver) RemoveAccount(args struct{ ID graphql.ID }) *accountOrFailResolver {
	defer r.crud.CloseCopy()
	id := bson.ObjectIdHex(string(args.ID))
	rawAccount, err := r.crud.FindOne(accountsCollection, &bson.M{"_id": id})
	genericErr := "Failed to remove account."
	if err != nil {
		return &accountOrFailResolver{&failResolver{
			"Invalid ID.",
		}}
	}
	account := transformAccount(rawAccount)
	err = r.crud.DeleteID(accountsCollection, id)
	if err != nil {
		fmt.Println("Failed to delete Account =>", err)
		return &accountOrFailResolver{&failResolver{genericErr}}
	}

	rawTokenMgr, err := r.crud.FindOne(tokenMgrCollection, &bson.M{"account_id": id})
	if err != nil {
		fmt.Println("Failed to find TokenManager =>", err)
		return &accountOrFailResolver{&failResolver{genericErr}}
	}

	tokenMgr := transformTokenManager(rawTokenMgr)
	err = r.crud.DeleteID(tokenMgrCollection, tokenMgr.ID)
	if err != nil {
		fmt.Println("Failed to delete TokenManager =>", err)
		return &accountOrFailResolver{&failResolver{genericErr}}
	}
	result := &accountResolver{&account}
	return &accountOrFailResolver{result}
}
