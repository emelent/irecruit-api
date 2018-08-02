package models

import (
	e "../errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Account db model
// look into excluding the password
type Account struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Email       string `json:"email" bson:"email"`
	Password    string `json:"-"  bson:"password"`
	Name        string `json:"name" bson:"name"`
	Surname     string `json:"surname" bson:"surname"`
	AccessLevel int    `json:"access_level" bson:"access_level"`

	HunterID  *bson.ObjectId `json:"hunter_id" bson:"hunter_id"`
	RecruitID *bson.ObjectId `json:"recruit_id" bson:"recruit_id"`
}

//OK validates Account fields
func (a *Account) OK() error {
	if a.Email == "" {
		return e.NewMissingFieldError("Email")
	}
	if a.Password == "" {
		return e.NewMissingFieldError("Password")
	}
	if a.Name == "" {
		return e.NewMissingFieldError("Name")
	}
	if a.Surname == "" {
		return e.NewMissingFieldError("Name")
	}

	return nil
}

// SetPassword sets account password
func (a *Account) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(bytes)
	return nil
}

// CheckPassword checks if given password matches the hash
func (a *Account) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// TODO implement
// func (a *Account) UpdatePassword(oldPassword, newPassword string) error {

// }
