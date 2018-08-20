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

// Login resolves graphql method of the same name
func (r *RootResolver) Login(ctx context.Context, args struct{ Email, Password string }) (*tokensResolver, error) {
	rawAccount, err := r.crud.FindOne(config.AccountsCollection, &bson.M{"email": args.Email})
	if err != nil {
		return nil, er.NewInvalidCredentialsError()
	}

	account := transformAccount(rawAccount)
	if !account.CheckPassword(args.Password) {
		return nil, er.NewInvalidCredentialsError()
	}

	rawTokenMgr, err := r.crud.FindOne(config.TokenManagersCollection, &bson.M{"account_id": account.ID})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.NewGenericError()
	}
	tokenMgr := transformTokenManager(rawTokenMgr)

	// create a new refresh token if the current one
	id := account.ID.Hex()
	ua := ctx.Value(mware.UaKey).(string)
	claims, err := utils.GetTokenClaims(tokenMgr.RefreshToken)
	if err != nil {
		log.Println("Invalid refresh token =>", err)
		return nil, er.NewGenericError()
	}
	t := time.Unix(claims.StandardClaims.ExpiresAt, 0)
	if time.Until(t) < time.Hour*24 {
		// create new refresh token
		tokenStr, err := utils.CreateRefreshToken(id)
		if err != nil {
			log.Println("Failed to create new refresh token =>", err)
			return nil, er.NewGenericError()
		}
		tokenMgr.RefreshToken = tokenStr

		// update refresh token in database
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

	return &tokensResolver{refresh: tokenMgr.RefreshToken, access: access}, nil

}
