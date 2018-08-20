package utils

import (
	"time"

	er "../errors"
	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("AllYourBase")

// Claims struct contains the jwt token claims
type Claims struct {
	UserAgent string `json:"user_agent"`
	AccountID string `json:"account_id"`
	Refresh   bool   `json:"refresh"`
	jwt.StandardClaims
}

func createToken(claims Claims) (*jwt.Token, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(mySigningKey)
	return token, signed, err
}

// CreateRefreshToken creates a refresh token
func CreateRefreshToken(accountID string) (string, error) {
	_, tokenStr, err := createToken(Claims{
		AccountID: accountID,
		Refresh:   true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24*30)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	return tokenStr, err
}

// CreateAccessToken creates an access token
func CreateAccessToken(accountID, ua string) (string, error) {
	_, tokenStr, err := createToken(Claims{
		AccountID: accountID,
		Refresh:   false,
		UserAgent: ua,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	return tokenStr, err
}

// GetTokenClaims parses the given token
func GetTokenClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, er.InvalidToken()

}
