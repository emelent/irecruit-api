package utils

import (
	e "../errors"
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

// CreateToken creates a jwt token with given claims
func CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(mySigningKey)
	return signed, err
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
	return nil, e.InvalidTokenError{}

}
