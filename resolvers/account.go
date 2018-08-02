package resolvers

import (
	"fmt"

	models "../models"
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
func (r *RootResolver) CreateAccount(args struct{ Info *accountDetails }) *accountResolver {
	defer r.crud.CloseCopy()
	account := models.Account{}
	info := args.Info
	account.Name = info.Name
	account.Email = info.Email
	account.Surname = info.Surname
	account.AccessLevel = 0
	account.SetPassword(info.Password)
	account.ID = bson.NewObjectId()

	err := r.crud.Insert(accountsCollection, account)
	if err != nil {
		fmt.Println("Error creating account =>", err)
		fmt.Printf("%#v", account)
	}
	return &accountResolver{&account}
}
