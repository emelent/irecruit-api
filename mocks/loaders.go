package mocks

import (
	"fmt"

	config "../config"
	db "../database"
	models "../models"
	utils "../utils"
)

// NewLoadedCRUD returns a crud object loaded with all the data
func NewLoadedCRUD() *db.CRUD {
	crud := db.NewCRUD(nil)

	// load data into db
	LoadAccounts(crud)
	LoadTokenManagers(crud)
	LoadRecruits(crud)
	LoadIndustries(crud)
	LoadQuestions(crud)
	LoadDocuments(crud)
	return crud
}

// LoadAccounts loads all account data
func LoadAccounts(crud *db.CRUD) {
	numRecruits := len(Recruits) - 1
	numHunters := len(HunterIDs) - 1

	for i, acc := range Accounts {
		acc.RecruitID = models.NullObjectID
		acc.HunterID = models.NullObjectID
		if i < numRecruits { // first n accounts have recruit profiles
			acc.RecruitID = Recruits[i].ID
		} else if i < numHunters+numRecruits { // next m accounts have hunter profiles
			acc.HunterID = HunterIDs[i-numRecruits]
		}

		if i == len(Accounts)-1 { // i.e. the SysAccount
			acc.HunterID = HunterIDs[numHunters]
			acc.RecruitID = Recruits[numRecruits].ID
		}
		//  set password to default
		acc.Password = DefaultPassword
		// validate before insertion
		if err := acc.OK(); err != nil {
			fmt.Printf("Mock accounts[%v] : %s", i, err.Error())
			break
		}

		// hash password before insertion
		acc.HashPassword()
		Accounts[i] = acc
		crud.Insert(config.AccountsCollection, acc)
	}
}

// LoadTokenManagers loads mock token managers
func LoadTokenManagers(crud *db.CRUD) {
	for i, mgr := range TokenManagers {
		mgr.AccountID = Accounts[i].ID

		//create refresh token
		refresh, err := utils.CreateRefreshToken(Accounts[i].ID.Hex())
		if err != nil {
			fmt.Printf("Mock token_manager[%v] : %s", i, err.Error())
			break
		}
		mgr.RefreshToken = refresh

		// validate before insertion
		if err := mgr.OK(); err != nil {
			fmt.Printf("Mock tokenManagers[%v] : %s", i, err.Error())
			break
		}

		TokenManagers[i] = mgr
		crud.Insert(config.TokenManagersCollection, mgr)
	}
}

// LoadRecruits loads mock recruits
func LoadRecruits(crud *db.CRUD) {
	for i, rec := range Recruits {
		// validate before insertion
		if err := rec.OK(); err != nil {
			fmt.Printf("Mock recruits[%v] : %s", i, err.Error())
			break
		}
		Recruits[i] = rec
		crud.Insert(config.RecruitsCollection, rec)
	}
}

// LoadIndustries loads mock industries
func LoadIndustries(crud *db.CRUD) {
	for i, industry := range Industries {
		// validate before insertion
		if err := industry.OK(); err != nil {
			fmt.Printf("Mock industries[%v] : %s", i, err.Error())
			break
		}

		Industries[i] = industry
		crud.Insert(config.IndustriesCollection, industry)
	}
}

// LoadQuestions load mock questions
func LoadQuestions(crud *db.CRUD) {
	var numIndustries = len(Industries)
	for i, q := range Questions {
		q.IndustryID = Industries[(i % numIndustries)].ID
		// validate before insertion
		if err := q.OK(); err != nil {
			fmt.Printf("Mock questions[%v] : %s", i, err.Error())
			break
		}

		Questions[i] = q
		crud.Insert(config.QuestionsCollection, q)
	}
}

// LoadDocuments load mock documents
func LoadDocuments(crud *db.CRUD) {
	var numRecruits = len(Recruits)
	for i, doc := range Documents {
		doc.OwnerID = Recruits[i%numRecruits].ID
		// validate before insertion
		if err := doc.OK(); err != nil {
			fmt.Printf("Mock documents[%v] : %s", i, err.Error())
			break
		}

		Documents[i] = doc
		crud.Insert(config.DocumentsCollection, doc)
	}
}
