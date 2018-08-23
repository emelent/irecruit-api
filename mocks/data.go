package mocks

import (
	models "../models"
	"gopkg.in/mgo.v2/bson"
)

// DefaultPassword is the password for all moc accounts
const DefaultPassword = "password"

/*
Contains all mock data

WARNING:

The data is arranged in a some-what specific way,
in order to make some logical assumptions about
the"connections" between the data within the loaders,
so moving things around could break tests and may
lead to unpredictable behaviour.

So, again, the structure of this data represents a specific
possible use case of the database.

DO NOT MOVE THE DATA AROUND UNLESS YOU UNDERSTAND
WHY IT IS WHERE IT IS.
*/

// HunterIDs 3 hunter IDs
var HunterIDs = []bson.ObjectId{
	bson.NewObjectId(),
	bson.NewObjectId(),
	bson.NewObjectId(),
	bson.NewObjectId(), // sys admin's hunterID
}

// Recruits 2 recruit profiles
var Recruits = []models.Recruit{
	{
		ID:         bson.NewObjectId(),
		Province:   "Gauteng",
		City:       "Pretoria",
		Gender:     "MALE",
		Disability: "",
		Vid1Url:    "none",
		Vid2Url:    "none",
		Phone:      "012 345 2378",
		Email:      "mark@gmail.com",
		Qa1:        models.QA{Question: "What's up?", Answer: "Nothing much."},
		Qa2:        models.QA{Question: "You good?", Answer: "You know it."},
		BirthYear:  1985,
	},
	{
		ID:         bson.NewObjectId(),
		Province:   "Gauteng",
		City:       "Johannesburg",
		Gender:     "MALE",
		Disability: "",
		Vid1Url:    "none",
		Vid2Url:    "none",
		Phone:      "013 345 2378",
		Email:      "johndoe@gmail.com",
		Qa1:        models.QA{Question: "What's in there?", Answer: "I don't know."},
		Qa2:        models.QA{Question: "Ever seen a turtle without it's shell?", Answer: "Nope."},
		BirthYear:  1995,
	},
	{ // sysadmin's recruitID
		ID:         bson.NewObjectId(),
		Province:   "North West",
		City:       "Mahikeng",
		Gender:     "MALE",
		Disability: "",
		Vid1Url:    "none",
		Vid2Url:    "none",
		Phone:      "014 345 2378",
		Email:      "thato@gmail.com",
		Qa1:        models.QA{Question: "What's this?", Answer: "A dead one of these."},
		Qa2:        models.QA{Question: "Ever seen a turtle without it's shell?", Answer: "All the time."},
		BirthYear:  1987,
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
		Name:        "Mark",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jdoe@gmail.com",
		Name:        "John",
		Surname:     "Doe",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "lisa@gmail.com",
		Name:        "Lisa",
		Surname:     "Smith",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "erin@gmail.com",
		Name:        "Erin",
		Surname:     "Lona",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "jake@gmail.com",
		Name:        "Jake",
		Surname:     "Tinder",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "moti@gmail.com",
		Name:        "Morlin",
		Surname:     "Tinder",
		AccessLevel: 0,
	},
	{
		ID:          bson.NewObjectId(),
		Email:       "thato@gmail.com",
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

// Documents 5 documents
var Documents = []models.Document{
	{ID: bson.NewObjectId(), URL: "http://google.com/d_1", DocType: "QUALIFICATION", OwnerType: "RECRUIT"},
	{ID: bson.NewObjectId(), URL: "http://google.com/d_2", DocType: "QUALIFICATION", OwnerType: "RECRUIT"},
	{ID: bson.NewObjectId(), URL: "http://google.com/d_3", DocType: "QUALIFICATION", OwnerType: "RECRUIT"},
	{ID: bson.NewObjectId(), URL: "http://google.com/d_4", DocType: "QUALIFICATION", OwnerType: "RECRUIT"},
	{ID: bson.NewObjectId(), URL: "http://google.com/d_5", DocType: "QUALIFICATION", OwnerType: "RECRUIT"},
}
