package mocks

import (
	"fmt"

	config "../config"
	db "../database"
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
	return crud
}

// LoadAccounts loads all account data
func LoadAccounts(crud *db.CRUD) {
	numRecruits := len(Recruits)
	numHunters := len(HunterIDs)

	for i, acc := range Accounts {
		if i < numRecruits { // first n accounts have recruit profiles
			acc.RecruitID = Recruits[i].ID
		} else if i < numHunters+numRecruits { // next m accounts have hunter profiles
			acc.HunterID = HunterIDs[i-numRecruits]
		}

		// validate before insertion
		if err := acc.OK(); err != nil {
			fmt.Printf("Mock accounts[%v] : %s", i, err.Error())
			break
		}

		// hash password before insertion
		acc.HashPassword()
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
		crud.Insert(config.IndustriesCollection, industry)
	}
}
