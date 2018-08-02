package resolvers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"context"

	mware "../middleware"
	utils "../utils"
	"gopkg.in/mgo.v2/bson"
)

const tokenMgrCollection = "token_managers"

// Login resolves graphql method of the same name
func (r *RootResolver) Login(ctx context.Context, args struct{ Email, Password string }) string {
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

	ua := ctx.Value(mware.UaKey).(string)
	// ip := ctx.Value("ip_address")

	//TODO create refresh token
	tokenStr, err := utils.CreateToken(utils.Claims{
		AccountID: account.ID.Hex(),
		Refresh:   true,
		UserAgent: ua,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24*30)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	if err != nil {
		fmt.Println("token error =>", err)
		return "Something went wrong."
	}

	//TODO register token in database
	// r.crud.Insert(tokenMgrCollection, Token{

	// })
	//TODO return refresh token
	return tokenStr
}
