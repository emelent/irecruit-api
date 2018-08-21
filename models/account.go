package models

import (
	"regexp"
	"strings"

	er "../errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Transformer
// -----------------

// TransformAccount transforms interface to Account model
func TransformAccount(in interface{}) Account {
	var account Account
	switch v := in.(type) {
	case bson.M:
		account.ID = v["_id"].(bson.ObjectId)
		account.Email = v["email"].(string)
		account.Name = v["name"].(string)
		account.Surname = v["surname"].(string)
		account.Password = v["password"].(string)
		account.AccessLevel = v["access_level"].(int)
		account.HunterID = v["hunter_id"].(bson.ObjectId)
		account.RecruitID = v["recruit_id"].(bson.ObjectId)
	case Account:
		account = v
	}

	return account
}

// -----------------
// Model
// -----------------

// Account db model
type Account struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Email       string `json:"email" bson:"email"`
	Password    string `json:"-"  bson:"password"`
	Name        string `json:"name" bson:"name"`
	Surname     string `json:"surname" bson:"surname"`
	AccessLevel int    `json:"access_level" bson:"access_level"`

	HunterID  bson.ObjectId `json:"hunter_id" bson:"hunter_id"`
	RecruitID bson.ObjectId `json:"recruit_id" bson:"recruit_id"`
}

//OK validates Account fields
func (a *Account) OK() error {
	reEmail := regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`)
	rePassword := regexp.MustCompile(`.{6}`)
	reName := regexp.MustCompile(`[a-zA-Z]{3,}`)
	reSurname := regexp.MustCompile(`[a-zA-Z]{3,}`)
	if !reEmail.MatchString(a.Email) {
		return er.InvalidField("Email")
	}
	if !rePassword.MatchString(a.Password) {
		return er.Input("Password must be at least 6 characters long.")
	}
	if !reName.MatchString(a.Name) {
		return er.Input("Name must be at least 3 alphabetic characters.")
	}
	if !reSurname.MatchString(a.Surname) {
		return er.Input("Surname must be at least 3 alphabetic characters.")
	}

	a.Email = strings.ToLower(a.Email)
	return nil
}

// HashPassword sets account password
func (a *Account) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
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
