package resolvers

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Login resolves graphql method of the same name
func (r *RootResolver) Login(args struct{ Email, Password string }) string {
	rawAccount, err := r.crud.FindOne(accountsCollection, bson.M{"email": args.Email})
	failedLogin := "Invalid username or email."
	if err != nil {
		fmt.Println("login error =>", err)
		return failedLogin
	}
	account := transformAccount(rawAccount)
	if !account.CheckPassword(args.Password) {
		return failedLogin
	}

	//TODO create refresh token
	//TODO register token in database
	//TODO return refresh token
	return "logged in"
}
