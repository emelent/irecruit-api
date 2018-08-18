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
		Qa1:        models.QA{Question: "What's up?", Answer: "Nothing much."},
		Qa2:        models.QA{Question: "You good?", Answer: "You know it."},
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
		Qa1:        models.QA{Question: "What's in there?", Answer: "I don't know."},
		Qa2:        models.QA{Question: "Ever seen a turtle without it's shell?", Answer: "Nope."},
	},
}

// TokenManagers 7  token managers
var TokenManagers = []models.TokenManager{
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
	{ID: bson.NewObjectId()},
}

// Accounts 7 user accounts
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
		Email:       "moti@gmail.com",
		Password:    "password",
		Name:        "Morlin",
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

// Industries 2 industries
var Industries = []models.Industry{
	{ID: bson.NewObjectId(), Name: "Statistics"},
	{ID: bson.NewObjectId(), Name: "Architecture"},
}

// Questions 5 questions
var Questions = []models.Question{
	// Industries[0] questions
	{ID: bson.NewObjectId(), Question: "What's your favourite colour?"},
	{ID: bson.NewObjectId(), Question: "What's your favourite song?"},
	{ID: bson.NewObjectId(), Question: "What's your favourite name?"},

	// Industries[1] questions
	{ID: bson.NewObjectId(), Question: "What's your favourite letter?"},
	{ID: bson.NewObjectId(), Question: "What's your favourite soup?"},
	{ID: bson.NewObjectId(), Question: "What's your favourite song?"},
}
