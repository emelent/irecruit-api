package resolvers

import (
	"context"
	"log"
	"time"

	config "../config"
	er "../errors"
	mware "../middleware"
	utils "../utils"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Root Resolver methods
// -----------------

// Login resolves graphql method "login"
func (r *RootResolver) Login(ctx context.Context, args struct{ Email, Password string }) (*TokensResolver, error) {
	defer r.crud.CloseCopy()

	// find account by email
	rawAccount, err := r.crud.FindOne(config.AccountsCollection, &bson.M{"email": args.Email})
	if err != nil {
		return nil, er.NewInvalidCredentialsError()
	}
	account := TransformAccount(rawAccount)

	// check if passwords match
	if !account.CheckPassword(args.Password) {
		return nil, er.NewInvalidCredentialsError()
	}

	// get account's tokenManager
	rawTokenMgr, err := r.crud.FindOne(config.TokenManagersCollection, &bson.M{"account_id": account.ID})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.NewGenericError()
	}
	tokenMgr := TransformTokenManager(rawTokenMgr)

	// get current refresh token
	claims, err := utils.GetTokenClaims(tokenMgr.RefreshToken)
	if err != nil {
		log.Println("Invalid refresh token =>", err)
		return nil, er.NewGenericError()
	}

	// prepare token creation data
	id := account.ID.Hex()
	ua := ctx.Value(mware.UaKey).(string)

	// create a new refresh token if the current one has expired
	t := time.Unix(claims.StandardClaims.ExpiresAt, 0)
	if time.Until(t) < time.Hour*24 {
		// create new refresh token
		tokenStr, err := utils.CreateRefreshToken(id)
		if err != nil {
			log.Println("Failed to create new refresh token =>", err)
			return nil, er.NewGenericError()
		}
		// update tokenManager's refresh token
		tokenMgr.RefreshToken = tokenStr

		// save tokenManager changes to db
		r.crud.UpdateID(
			config.TokenManagersCollection,
			tokenMgr.ID, bson.M{"refresh_token": tokenStr},
		)
	}

	// create access token
	access, err := utils.CreateAccessToken(id, ua)
	if err != nil {
		log.Println("Failed to create access token =>", err)
		return nil, er.NewGenericError()
	}

	return &TokensResolver{refresh: tokenMgr.RefreshToken, access: access}, nil

}
