package resolvers

import (
	"log"
	"time"

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

	// get token claims
	claims, err := utils.GetTokenClaims(args.Token)
	if err != nil {
		return nil, er.InvalidToken()
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
	account := models.TransformAccount(rawAccount)

	// func to resolve Viewer as RecruitViewer
	viewAsRecruit := func() (*ViewerResolver, error) {
		// check if account has recruit profile
		if utils.IsNullID(account.RecruitID) {
			return nil, er.Input("Failed to enforce 'RECRUIT'.")
		}

		// retrieve Recruit profile
		rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
		if err != nil {
			log.Println("Failed to find recruit =>", err)
			return nil, er.Generic()
		}

		// return RecruitViewer
		recruit := models.TransformRecruit(rawRecruit)
		viewer := &RecruitViewerResolver{&recruit, &account, r.crud}
		return &ViewerResolver{viewer}, nil
	}

	// func to resolve Viewer as SysViewer
	viewAsSys := func() (*ViewerResolver, error) {
		// check if account is sys account
		if !utils.IsSysAccount(&account) {
			return nil, er.Input("Failed to enforce 'SYSTEM'.")
		}

		// return sysViewer
		viewer := &SysViewerResolver{&account, r.crud}
		return &ViewerResolver{viewer}, nil
	}

	//func to resolve Viewer as AccountViewer
	viewAsAccount := func() (*ViewerResolver, error) {
		// return accountViewer
		viewer := &AccountViewerResolver{&account}
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
		return viewer, nil
	}

	// try to view as RecruitViewer
	if viewer, err := viewAsRecruit(); err == nil {
		return viewer, nil
	}

	// if all else fails view as AccountEditor
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
	account := models.TransformAccount(rawAccount)

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
// 	recruit  := models.TransformRecruit(rawRecruit)
// 	rawAccount, err := r.crud.FindOne(
// 		config.AccountsCollection,
// 		&bson.M{"recruit_id": recruit.ID},
// 	)
// 	if err !=  nil{
// 		log.Println(err)
// 		return nil, er.GenericError()
// 	}
// 	account :=  models.TransformAccount(rawAccount)
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
		account := models.TransformAccount(r)
		results = append(results, &AccountResolver{&account})
	}
	return results, err
}

// Recruits resolves SysViewer.Recruits which returns a list of all the recruits
func (r *SysViewerResolver) Recruits() ([]*RecruitResolver, error) {
	defer r.crud.CloseCopy()

	// fetch all recruits
	rawRecruits, err := r.crud.FindAll(config.RecruitsCollection, nil)
	if err != nil {
		return nil, er.Generic()
	}

	// process results
	results := make([]*RecruitResolver, 0)
	for _, raw := range rawRecruits {
		recruit := models.TransformRecruit(raw)
		results = append(results, &RecruitResolver{&recruit, r.a})
	}

	// return results
	return results, nil
}

// Questions resolves SysViewer.Questions which returns a list of all the questions
func (r *SysViewerResolver) Questions() ([]*QuestionResolver, error) {
	defer r.crud.CloseCopy()

	// fetch all recruits
	rawQuestions, err := r.crud.FindAll(config.QuestionsCollection, nil)
	if err != nil {
		return nil, er.Generic()
	}

	// process results
	results := make([]*QuestionResolver, 0)
	for _, raw := range rawQuestions {
		question := models.TransformQuestion(raw)
		results = append(results, &QuestionResolver{&question})
	}

	// return results
	return results, nil
}

// Documents resolves SysViewer.Documents which returns a list of all the documents
func (r *SysViewerResolver) Documents() ([]*DocumentResolver, error) {
	defer r.crud.CloseCopy()

	// fetch all documents
	rawDocuments, err := r.crud.FindAll(config.DocumentsCollection, nil)
	if err != nil {
		return nil, er.Generic()
	}

	// process results
	results := make([]*DocumentResolver, 0)
	for _, raw := range rawDocuments {
		document := models.TransformDocument(raw)
		results = append(results, &DocumentResolver{&document})
	}

	// return results
	return results, nil
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

// -----------------
// QaResolver struct
// -----------------

// QaResolver resolve Qa
type QaResolver struct {
	qa *models.QA
}

// Question resolves Qa.Question
func (r *QaResolver) Question() string {
	return r.qa.Question
}

// Answer resolves Qa.Answer
func (r *QaResolver) Answer() string {
	return r.qa.Answer
}

// -----------------
// RecruitResolver struct
// -----------------

// RecruitResolver resolves Recruit
type RecruitResolver struct {
	r *models.Recruit
	a *models.Account
}

// ID resolves Recruit.ID
func (r *RecruitResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
}

// Age resolves Recruit.Age
func (r *RecruitResolver) Age() int32 {
	year := int32(time.Now().Year())
	return year - r.r.BirthYear
}

// Name resolves Recruit.Name
func (r *RecruitResolver) Name() string {
	return r.a.Name
}

// Surname resolves Recruit.Surname
func (r *RecruitResolver) Surname() string {
	return r.a.Surname
}

// Phone resolves Recruit.Phone
func (r *RecruitResolver) Phone() string {
	return r.r.Phone
}

// Email resolves Recruit.Email
func (r *RecruitResolver) Email() string {
	return r.r.Email
}

// Province resolves Recruit.Province
func (r *RecruitResolver) Province() string {
	return r.r.Province
}

// City resolves Recruit.City
func (r *RecruitResolver) City() string {
	return r.r.City
}

// Gender resolves Recruit.Gender
func (r *RecruitResolver) Gender() string {
	return r.r.Gender
}

// Disability resolves Recruit.Disability
func (r *RecruitResolver) Disability() string {
	return r.r.Disability
}

// Vid1Url resolves Recruit.Vid1Url
func (r *RecruitResolver) Vid1Url() string {
	return r.r.Vid1Url
}

// Vid2Url resolves Recruit.Vid2Url
func (r *RecruitResolver) Vid2Url() string {
	return r.r.Vid2Url
}

// Qa1 resolves Recruit.Qa1
func (r *RecruitResolver) Qa1() *QaResolver {
	return &QaResolver{&r.r.Qa1}
}

// Qa2 resolves Recruit.Qa2
func (r *RecruitResolver) Qa2() *QaResolver {
	return &QaResolver{&r.r.Qa2}
}
