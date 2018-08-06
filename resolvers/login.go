package resolvers

import (
	"context"
	"fmt"
	"time"

	mware "../middleware"
	utils "../utils"
	"gopkg.in/mgo.v2/bson"
)

const tokenMgrCollection = "token_managers"

// Login resolves graphql method of the same name
func (r *RootResolver) Login(ctx context.Context, args struct{ Email, Password string }) *tokensOrFailResolver {
	rawAccount, err := r.crud.FindOne(accountsCollection, &bson.M{"email": args.Email})
	failedLogin := "Invalid username or email."
	genericErr := "Something went wrong."
	if err != nil {
		return &tokensOrFailResolver{&failResolver{failedLogin}}
	}

	account := transformAccount(rawAccount)
	if !account.CheckPassword(args.Password) {
		return &tokensOrFailResolver{&failResolver{failedLogin}}
	}

	rawTokenMgr, err := r.crud.FindOne(tokenMgrCollection, &bson.M{"account_id": account.ID})
	if err != nil {
		fmt.Println("Failed to find TokenManager =>", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
	}
	tokenMgr := transformTokenManager(rawTokenMgr)

	// create a new refresh token if the current one
	id := account.ID.Hex()
	ua := ctx.Value(mware.UaKey).(string)
	claims, err := utils.GetTokenClaims(tokenMgr.RefreshToken)
	t := time.Unix(claims.StandardClaims.ExpiresAt, 0)
	if time.Until(t) < time.Hour*24 {
		// create new refresh token
		tokenStr, err := utils.CreateRefreshToken(id)
		if err != nil {
			fmt.Println("Failed to create new refresh token =>", err)
			return &tokensOrFailResolver{&failResolver{genericErr}}
		}
		tokenMgr.RefreshToken = tokenStr

		// update refresh token in database
		r.crud.UpdateID(
			tokenMgrCollection,
			tokenMgr.ID, bson.M{"refresh_token": tokenStr},
		)
	}

	// create access token
	access, err := utils.CreateAccessToken(id, ua)
	if err != nil {
		fmt.Println("Failed to create access token =>", err)
		return &tokensOrFailResolver{&failResolver{genericErr}}
	}

	result := &tokensResolver{refresh: tokenMgr.RefreshToken, access: access}
	return &tokensOrFailResolver{result}
}
