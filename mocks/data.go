package mocks

import (
	models "../models"
	"gopkg.in/mgo.v2/bson"
)

// Contains all mock data

// HunterIDs 3 hunter IDs
var HunterIDs = []bson.ObjectId{
	bson.NewObjectId(),
	bson.NewObjectId(),
	bson.NewObjectId(),
}

// Recruits 2 recruit profiles
var Recruits = []models.Recruit{
	{
		ID:         bson.NewObjectId(),
		Province:   "Gauteng",
		City:       "Pretoria",
		Gender:     "male",
		Disability: "",
		Vid1Url:    "none",
		Vid2Url:    "none",
		Phone:      "012 345 2378",
		Email:      "mark@gmail.com",
	},
	{
		ID:         bson.NewObjectId(),
		Province:   "Gauteng",
		City:       "Johannesburg",
		Gender:     "male",
		Disability: "",
		Vid1Url:    "none",
		Vid2Url:    "none",
		Phone:      "013 345 2378",
		Email:      "johndoe@gmail.com",
	},
}

// TokenManagers 6  token managers
var TokenManagers = []models.TokenManager{
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
}

// Accounts 6 user accounts
var Accounts = []models.Account{
	{
		ID:          bson.NewObjectId(),
		Email:       "mark@gmail.com",
		Password:    "password",
		Name:        "Mark",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jdoe@gmail.com",
		Password:    "password",
		Name:        "John",
		Surname:     "Doe",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "lisa@gmail.com",
		Password:    "password",
		Name:        "Lisa",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "erin@gmail.com",
		Password:    "password",
		Name:        "Erin",
		Surname:     "Lona",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jake@gmail.com",
		Password:    "password",
		Name:        "Jake",
		Surname:     "Tinder",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "thato@gmail.com",
		Password:    "password",
		Name:        "Thato",
		Surname:     "Mopani",
		AccessLevel: 9, // system admin
	},
}
