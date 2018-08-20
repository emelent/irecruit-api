package resolvers

import (
	"log"

	config "../config"
	db "../database"
	er "../errors"
	models "../models"
	utils "../utils"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Root Resolver Methods struct
// -----------------

// View resolves "view" gql query
func (r *RootResolver) View(args struct {
	Token   string
	Enforce *string
}) (*ViewerResolver, error) {

	// Did we get a token?
	if args.Token == "" {
		// return a guest viewer
		return nil, er.MissingField("token")
	}

	// get token claims
	claims, err := utils.GetTokenClaims(args.Token)
	if err != nil {
		return nil, err
	}

	// no refresh tokens allowed
	if claims.Refresh {
		return nil, er.InvalidToken()
	}

	// check claims AccountID
	if !bson.IsObjectIdHex(claims.AccountID) {
		return nil, er.InvalidToken()
	}

	// get account
	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(claims.AccountID))
	if err != nil {
		log.Println("Failed to find account by ID from token =>", err)
		return nil, er.InvalidToken()
	}
	account := TransformAccount(rawAccount)

	// func to resolve Viewer as RecruitViewer
	viewAsRecruit := func() (*ViewerResolver, error) {
		// check if account has recruit profile
		if utils.IsNullID(account.RecruitID) {
			return nil, er.Input("Failed to enfore 'RECRUIT'.")
		}

		// retrieve Recruit profile
		rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
		if err != nil {
			log.Println("Failed to find recruit =>", err)
			return nil, er.Generic()
		}

		// return RecruitViewer
		recruit := TransformRecruit(rawRecruit)
		viewer := &RecruitViewerResolver{&recruit, &account, r.crud}
		return &ViewerResolver{viewer}, nil
	}

	// func to resolve Viewer as SysViewer
	viewAsSys := func() (*ViewerResolver, error) {
		// check if account is sys account
		if !utils.IsSysAccount(&account) {
			return nil, er.Input("Failed to enfore 'SYSTEM'.")
		}

		// return sysViewer
		viewer := &SysViewerResolver{&account, r.crud}
		return &ViewerResolver{viewer}, nil
	}

	//func to resolve Viewer as AccountViewer
	viewAsAccount := func() (*ViewerResolver, error) {
		// check if account is sys account
		if !utils.IsSysAccount(&account) {
			return nil, er.Input("Failed to enfore 'SYSTEM'.")
		}

		// return sysViewer
		viewer := &SysViewerResolver{&account, r.crud}
		return &ViewerResolver{viewer}, nil
	}

	// enforce a Viewer type if specified
	if args.Enforce != nil {
		switch *args.Enforce {
		case "RECRUIT":
			return viewAsRecruit()
		case "HUNTER":
			return nil, er.Input("Unimplemented")
		case "SYSTEM":
			return viewAsSys()
		case "ACCOUNT":
			return viewAsAccount()
		default:
			return nil, er.InvalidField("enforce")
		}
	}

	// try to view as SysViewer
	if viewer, err := viewAsSys(); err == nil {
		return viewer, err
	}

	// try to view as RecruitViewer
	if viewer, err := viewAsRecruit(); err == nil {
		return viewer, err
	}

	// view as just a bare account
	return viewAsAccount()
}

// -----------------
// viewer interface
// -----------------
type viewer interface {
	ID() graphql.ID
	Name() string
	Surname() string
	Email() string
}

// -----------------
// RecruitViewerResolver struct
// -----------------

// RecruitViewerResolver resolves RecruitViewer
type RecruitViewerResolver struct {
	r    *models.Recruit
	a    *models.Account
	crud *db.CRUD
}

// ID resolves RecruitViewer.ID
func (r *RecruitViewerResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
}

// Name resolves RecruitViewer.Name
func (r *RecruitViewerResolver) Name() string {
	return r.a.Name
}

// Surname resolves RecruitViewer.Surname
func (r *RecruitViewerResolver) Surname() string {
	return r.a.Surname
}

// Email resolves RecruitViewer.Email
func (r *RecruitViewerResolver) Email() string {
	return r.a.Email
}

// Profile resolves RecruitViewer.Profile which returns the current account's Recruit profile
func (r *RecruitViewerResolver) Profile() (*RecruitResolver, error) {
	defer r.crud.CloseCopy()

	// retrieve account
	rawAccount, err := r.crud.FindID(config.AccountsCollection, r.a.ID)
	if err != nil {
		log.Println(err)
		return nil, er.Generic()
	}
	account := TransformAccount(rawAccount)

	return &RecruitResolver{r.r, &account}, nil
}

// -----------------
// hunterViewerResolver struct
// -----------------
// type hunterViewerResolver struct {
// 	r *models.Hunter
// 	crud *db.CRUD
// }

// func (r *hunterViewerResolver) ID() graphql.ID {
// 	return graphql.ID(r.r.ID.Hex())
// }

// func (r *hunterViewerResolver) Name() string {
// 	return graphql.ID(r.r.Name)
// }

// func (r *hunterViewerResolver) Surname() string {
// 	return graphql.ID(r.r.Surname)
// }

// func (r *hunterViewerResolver) Email() string {
// 	return graphql.ID(r.r.Email)
// }

// func (r *hunterViewerResolver) Recruit(args struct{ID graphql.ID}) (*RecruitResolver, error){
// 	id := string(args.ID)
// 	if !bson.IsObjectIdHex(id){
// 		log.Println("Invalid id!")
// 		return nil, er.GenericError()
// 	}

// 	rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
// 	if err !=  nil{
// 		log.Println(err)
// 		return nil, er.GenericError()
// 	}
// 	recruit  := TransformRecruit(rawRecruit)
// 	rawAccount, err := r.crud.FindOne(
// 		config.AccountsCollection,
// 		&bson.M{"recruit_id": recruit.ID},
// 	)
// 	if err !=  nil{
// 		log.Println(err)
// 		return nil, er.GenericError()
// 	}
// 	account :=  TransformAccount(rawAccount)
// 	return &RecruitResolver{&recruit, &account}, nil
// }

// -----------------
// SysViewerResolver struct
// -----------------

// SysViewerResolver resolves SysViewer
type SysViewerResolver struct {
	a    *models.Account
	crud *db.CRUD
}

// ID resolves SysViewer.ID
func (r *SysViewerResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

//Name resolves SysViewer.Name
func (r *SysViewerResolver) Name() string {
	return r.a.Name
}

// Surname resolves SysViewer.Surname
func (r *SysViewerResolver) Surname() string {
	return r.a.Surname
}

// Email resolves SysViewer.Email
func (r *SysViewerResolver) Email() string {
	return r.a.Email
}

// Accounts resolves SysViewer.Accounts which returns a list of all the accounts
func (r *SysViewerResolver) Accounts() ([]*AccountResolver, error) {
	defer r.crud.CloseCopy()

	// fetch all accounts
	rawAccounts, err := r.crud.FindAll(config.AccountsCollection, nil)
	results := make([]*AccountResolver, 0)
	for _, r := range rawAccounts {
		account := TransformAccount(r)
		results = append(results, &AccountResolver{&account})
	}
	return results, err
}

// -----------------
// AccountViewerResolver struct
// -----------------

// AccountViewerResolver resolves AccountViewer
type AccountViewerResolver struct {
	a *models.Account
}

// ID resolves AccountViewer.ID
func (r *AccountViewerResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

// Name resolves AccountViewer.Name
func (r *AccountViewerResolver) Name() string {
	return r.a.Name
}

// Surname resolves AccountViewer.Surname
func (r *AccountViewerResolver) Surname() string {
	return r.a.Surname
}

// Email resolves AccountViewer.Email
func (r *AccountViewerResolver) Email() string {
	return r.a.Email
}

// IsHunter resolves AccountViewer.IsHunter
func (r *AccountViewerResolver) IsHunter() bool {
	return !utils.IsNullID(r.a.HunterID)
}

// IsRecruit resolves AccountViewer.IsRecruit
func (r *AccountViewerResolver) IsRecruit() bool {
	return !utils.IsNullID(r.a.RecruitID)
}

// CheckPassword resolves AccountViewer.CheckPassword
func (r *AccountViewerResolver) CheckPassword(args struct{ Password string }) bool {
	return r.a.CheckPassword(args.Password)
}

// -----------------
// ViewerResolver struct
// -----------------

// ViewerResolver resolves Viewer
type ViewerResolver struct {
	viewer
}

// ToRecruitViewer asserts *ViewerResolver to *RecruitViewerResolver
func (r *ViewerResolver) ToRecruitViewer() (*RecruitViewerResolver, bool) {
	v, ok := r.viewer.(*RecruitViewerResolver)
	return v, ok
}

// ToSysViewer asserts *ViewerResolver to *SysViewerResolver
func (r *ViewerResolver) ToSysViewer() (*SysViewerResolver, bool) {
	v, ok := r.viewer.(*SysViewerResolver)
	return v, ok
}

// ToAccountViewer asserts *ViewerResolver to *AccountViewerResolver
func (r *ViewerResolver) ToAccountViewer() (*AccountViewerResolver, bool) {
	v, ok := r.viewer.(*AccountViewerResolver)
	return v, ok
}
