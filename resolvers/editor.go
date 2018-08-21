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

// Edit resolves "view" gql query
func (r *RootResolver) Edit(args struct {
	Token   string
	Enforce *string
}) (*EditorResolver, error) {

	// Did we get a token?
	if args.Token == "" {
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

	// get account
	if !bson.IsObjectIdHex(claims.AccountID) {
		return nil, er.InvalidToken()
	}

	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(claims.AccountID))
	if err != nil {
		log.Println("Failed to find account by ID from token =>", err)
		return nil, er.InvalidToken()
	}
	account := models.TransformAccount(rawAccount)

	// enforce an enforceable
	if args.Enforce != nil {
		switch *args.Enforce {
		case "RECRUIT": // try to enforce recruit

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

			// return RecruitEditor
			recruit := models.TransformRecruit(rawRecruit)
			Editor := &RecruitEditorResolver{&recruit, &account, r.crud}
			return &EditorResolver{Editor}, nil

		case "HUNTER": // try to enforce hunter
			return nil, er.Input("Unimplemented")

		case "SYSTEM": // try to enforce system
			// check if account is sys account
			if !utils.IsSysAccount(&account) {
				return nil, er.Input("Failed to enfore 'SYSTEM'.")
			}

			// return sysEditor
			Editor := &SysEditorResolver{&account, r.crud}
			return &EditorResolver{Editor}, nil

		case "ACCOUNT":
			// return accountEditor
			return &EditorResolver{&AccountEditorResolver{&account, r.crud}}, nil

		}
	}

	// check if account is sys account
	if utils.IsSysAccount(&account) {
		Editor := &SysEditorResolver{&account, r.crud}
		return &EditorResolver{Editor}, nil
	}

	// check if account has recruit profile
	if !utils.IsNullID(account.RecruitID) {
		rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
		if err != nil {
			log.Println("Failed to find recruit =>", err)
			return nil, er.Generic()
		}
		recruit := models.TransformRecruit(rawRecruit)
		Editor := &RecruitEditorResolver{&recruit, &account, r.crud}
		return &EditorResolver{Editor}, nil
	}

	return nil, er.Generic()
}

// -----------------
// editor interface
// -----------------

type editor interface{}

// -----------------
// RecruitEditorResolver struct
// -----------------

// RecruitEditorResolver resolves RecruitEditor
type RecruitEditorResolver struct {
	r    *models.Recruit
	a    *models.Account
	crud *db.CRUD
}

// RemoveRecruit resolves "removeRecruit" mutation
func (r *RecruitEditorResolver) RemoveRecruit() (*string, error) {
	defer r.crud.CloseCopy()

	// attempt to remove Recruit
	if err := r.crud.DeleteID(config.RecruitsCollection, r.r.ID); err != nil {
		return nil, er.Generic()
	}

	// remove the account's recruit_id
	if err := r.crud.UpdateID(config.AccountsCollection, r.a.ID, bson.M{
		"recruit_id": models.NullObjectID,
	}); err != nil {
		return nil, er.Generic()
	}

	result := "Recruit successfully removed."
	return &result, nil
}

// -----------------
// hunterEditorResolver struct
// -----------------

// type hunterEditorResolver struct {
// 	r *models.Hunter
// 	crud *db.CRUD
// }

// -----------------
// SysEditorResolver struct
// -----------------

// SysEditorResolver resolves SysEditor
type SysEditorResolver struct {
	a    *models.Account
	crud *db.CRUD
}

// ID resolves SysEditor.ID
func (r *SysEditorResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

// Name resolves SysEditor.Name
func (r *SysEditorResolver) Name() string {
	return r.a.Name
}

// Surname resolves SysEditor.Surname
func (r *SysEditorResolver) Surname() string {
	return r.a.Surname
}

// Email resolves SysEditor.Email
func (r *SysEditorResolver) Email() string {
	return r.a.Email
}

// RemoveRecruit resolves SysEditor.RemoveRecruit which removes a Recruit with the given ID
func (r *SysEditorResolver) RemoveRecruit(args struct{ ID graphql.ID }) (*string, error) {
	return ResolveRemoveByID(
		r.crud,
		config.RecruitsCollection,
		"Recruit",
		string(args.ID),
	)
}

// RemoveAccount resolves SysEditor.RemoveAccount which removes an Account with the given ID
func (r *SysEditorResolver) RemoveAccount(args struct{ ID graphql.ID }) (*string, error) {
	id := string(args.ID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.InvalidField("id")
	}

	return ResolveRemoveAccount(r.crud, bson.ObjectIdHex(id))
}

// -----------------
// AccountEditorResolver struct
// -----------------

// AccountEditorResolver resolves AccountEditor
type AccountEditorResolver struct {
	a    *models.Account
	crud *db.CRUD
}

// RemoveAccount resolves AccountEditor.RemoveAccount which removes the current account
func (r *AccountEditorResolver) RemoveAccount() (*string, error) {
	return ResolveRemoveAccount(r.crud, r.a.ID)
}

// CreateRecruit resolves AccountEditor.CreateRecruit which creates a Recruit profile for the current account using the given Info
func (r *AccountEditorResolver) CreateRecruit(args struct{ Info *recruitDetails }) (*RecruitResolver, error) {
	// check if the account has a recruit profile
	account := r.a
	if !utils.IsNullID(account.RecruitID) {
		return nil, er.Input("Account already has a Recruit profile.")
	}

	// check if info is nil
	info := args.Info
	if info == nil {
		return nil, er.MissingField("info")
	}

	// validate info
	if info.Province == nil {
		return nil, er.MissingField("info.province")
	}
	if info.Phone == nil {
		return nil, er.MissingField("info.phone")
	}
	if info.Email == nil {
		return nil, er.MissingField("info.email")
	}
	if info.City == nil {
		return nil, er.MissingField("info.city")
	}
	if info.Gender == nil {
		return nil, er.MissingField("info.gender")
	}
	if info.Disability == nil {
		return nil, er.MissingField("info.disability")
	}
	if info.Vid1Url == nil {
		return nil, er.MissingField("info.vid1_url")
	}
	if info.Vid2Url == nil {
		return nil, er.MissingField("info.vid2_url")
	}
	if info.BirthYear == nil {
		return nil, er.MissingField("info.birth_year")
	}
	if info.Qa1Question == nil {
		return nil, er.MissingField("info.qa1_question")
	}
	if info.Qa1Answer == nil {
		return nil, er.MissingField("info.qa1_answer")
	}
	if info.Qa2Question == nil {
		return nil, er.MissingField("info.qa2_question")
	}
	if info.Qa2Answer == nil {
		return nil, er.MissingField("info.qa2_answer")
	}

	// create recruit profile
	var recruit models.Recruit
	recruit.ID = bson.NewObjectId()
	recruit.Province = *info.Province
	recruit.City = *info.City
	recruit.Gender = *info.Gender
	recruit.Disability = *info.Disability
	recruit.Vid1Url = *info.Vid1Url
	recruit.Vid2Url = *info.Vid2Url
	recruit.Phone = *info.Phone
	recruit.Email = *info.Email
	recruit.BirthYear = *info.BirthYear
	recruit.Qa1 = models.QA{Question: *info.Qa1Question, Answer: *info.Qa1Answer}
	recruit.Qa2 = models.QA{Question: *info.Qa2Question, Answer: *info.Qa2Answer}

	// validate recruit profile
	if err := recruit.OK(); err != nil {
		return nil, err
	}

	// store recruit profile in database
	if err := r.crud.Insert(config.RecruitsCollection, recruit); err != nil {
		log.Println(err)
		return nil, er.Generic()
	}

	// attach the recruit profile to the account
	if err := r.crud.UpdateID(config.AccountsCollection, account.ID, bson.M{
		"recruit_id": recruit.ID,
	}); err != nil {
		log.Println(err)
		return nil, er.Generic()
	}
	return &RecruitResolver{&recruit, account}, nil
}

// -----------------
// EditorResolver struct
// -----------------

// EditorResolver resolves Editor
type EditorResolver struct {
	editor
}

// ToRecruitEditor asserts *EditorResolver to *RecruitEditorResolver
func (r *EditorResolver) ToRecruitEditor() (*RecruitEditorResolver, bool) {
	v, ok := r.editor.(*RecruitEditorResolver)
	return v, ok
}

// ToSysEditor asserts *EditorResolver to *SysEditorResolver
func (r *EditorResolver) ToSysEditor() (*SysEditorResolver, bool) {
	v, ok := r.editor.(*SysEditorResolver)
	return v, ok
}

// ToAccountEditor asserts *EditorResolver to *AccountEditorResolver
func (r *EditorResolver) ToAccountEditor() (*AccountEditorResolver, bool) {
	v, ok := r.editor.(*AccountEditorResolver)
	return v, ok
}

// -----------------
// recruitDetails struct
// -----------------
type recruitDetails struct {
	Phone       *string
	Email       *string
	Province    *string
	City        *string
	Gender      *string
	Disability  *string
	Vid1Url     *string
	Vid2Url     *string
	Qa1Question *string
	Qa1Answer   *string
	Qa2Question *string
	Qa2Answer   *string
	BirthYear   *int32
}
