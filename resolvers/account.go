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

type accountResolver struct {
	a *models.Account
}

type accountDetails struct {
	Email    string
	Password string
	Name     string
	Surname  string
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
func (r *RootResolver) Accounts(args struct{ Name *string }) []*accountResolver {
	defer r.crud.CloseCopy()
	if args.Name == nil {
		// TODO return all names
	}
	rawAccounts, _ := r.crud.FindAll(accountsCollection, nil)
	results := make([]*accountResolver, 0)
	for _, r := range rawAccounts {
		account := transformAccount(r)
		results = append(results, &accountResolver{&account})
	}
	return results
}

// CreateAccount resolves the query of the same name
func (r *RootResolver) CreateAccount(ctx context.Context, args struct{ Info *accountDetails }) *tokensResolver {
	defer r.crud.CloseCopy()
	account := models.Account{}
	info := args.Info
	account.Name = info.Name
	account.Email = info.Email
	account.Surname = info.Surname
	account.AccessLevel = 0
	account.SetPassword(info.Password)
	account.ID = bson.NewObjectId()

	// create account
	err := r.crud.Insert(accountsCollection, account)
	if err != nil {
		fmt.Println("Error creating account =>", err)
		fmt.Printf("%#v", account)
		return nil
	}
	id := account.ID.Hex()
	// create refresh token
	refresh, err := utils.CreateRefreshToken(id)

	if err != nil {
		fmt.Println("Error creating refresh token =>", err)
		return nil
	}

	// access token
	ua := ctx.Value(mware.UaKey).(string)
	access, err := utils.CreateAccessToken(id, ua)

	if err != nil {
		fmt.Println("Error creating access token =>", err)
		return nil
	}

	// create token manager
	tokenMgr := models.TokenManager{
		AccountID:    account.ID,
		Tokens:       []string{access},
		RefreshToken: refresh,
		MaxTokens:    5,
	}
	err = r.crud.Insert(tokenMgrCollection, tokenMgr)
	if err != nil {
		fmt.Println("Error creating token manager =>", err)
		fmt.Printf("%#v", account)
		return nil
	}

	return &tokensResolver{refresh: refresh, access: access}
}

// RemoveAccount removes an account
func (r *RootResolver) RemoveAccount(args struct{ ID graphql.ID }) string {
	defer r.crud.CloseCopy()

	id := bson.ObjectIdHex(string(args.ID))
	err := r.crud.DeleteID(accountsCollection, id)
	if err != nil {
		return "Failed to remove account."
	}
	return "Account successfully removed."
}
