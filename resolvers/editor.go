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
}) (*editorResolver, error) {

	// Did we get a token?
	if args.Token == "" {
		return nil, er.NewMissingFieldError("token")
	}

	// get token claims
	claims, err := utils.GetTokenClaims(args.Token)
	if err != nil {
		return nil, err
	}

	// no refresh tokens allowed
	if claims.Refresh {
		return nil, er.NewInvalidTokenError()
	}

	// get account
	if !bson.IsObjectIdHex(claims.AccountID) {
		return nil, er.NewInvalidTokenError()
	}

	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(claims.AccountID))
	if err != nil {
		log.Println("Failed to find account by ID from token =>", err)
		return nil, er.NewInvalidTokenError()
	}
	account := transformAccount(rawAccount)

	// enforce an enforceable
	if args.Enforce != nil {
		switch *args.Enforce {
		case "RECRUIT": // try to enforce recruit

			// check if account has recruit profile
			if utils.IsNullID(account.RecruitID) {
				return nil, er.NewInputError("Failed to enfore 'RECRUIT'.")
			}

			// retrieve Recruit profile
			rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
			if err != nil {
				log.Println("Failed to find recruit =>", err)
				return nil, er.NewGenericError()
			}

			// return RecruitEditor
			recruit := transformRecruit(rawRecruit)
			editor := &recruitEditorResolver{&recruit, &account, r.crud}
			return &editorResolver{editor}, nil

		case "HUNTER": // try to enforce hunter
			return nil, er.NewInputError("Unimplemented")

		case "SYSTEM": // try to enforce system
			// check if account is sys account
			if !utils.IsSysAccount(&account) {
				return nil, er.NewInputError("Failed to enfore 'SYSTEM'.")
			}

			// return sysEditor
			editor := &sysEditorResolver{&account, r.crud}
			return &editorResolver{editor}, nil

		case "ACCOUNT":
			// return accountEditor
			return &editorResolver{&accountEditorResolver{&account, r.crud}}, nil

		}
	}

	// check if account is sys account
	if utils.IsSysAccount(&account) {
		editor := &sysEditorResolver{&account, r.crud}
		return &editorResolver{editor}, nil
	}

	// check if account has recruit profile
	if !utils.IsNullID(account.RecruitID) {
		rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
		if err != nil {
			log.Println("Failed to find recruit =>", err)
			return nil, er.NewGenericError()
		}
		recruit := transformRecruit(rawRecruit)
		editor := &recruitEditorResolver{&recruit, &account, r.crud}
		return &editorResolver{editor}, nil
	}

	return nil, er.NewGenericError()
}

// -----------------
// editor interface
// -----------------
type editor interface{}

// -----------------
// recruitEditorResolver struct
// -----------------
type recruitEditorResolver struct {
	r    *models.Recruit
	a    *models.Account
	crud *db.CRUD
}

// RemoveRecruit resolves "removeRecruit" mutation
func (r *recruitEditorResolver) RemoveRecruit() (*string, error) {
	defer r.crud.CloseCopy()

	// attempt to remove Recruit
	if err := r.crud.DeleteID(config.RecruitsCollection, r.r.ID); err != nil {
		return nil, er.NewGenericError()
	}

	// remove the account's recruit_id
	if err := r.crud.UpdateID(config.AccountsCollection, r.a.ID, bson.M{
		"recruit_id": models.NullObjectID,
	}); err != nil {
		return nil, er.NewGenericError()
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
// sysEditorResolver struct
// -----------------
type sysEditorResolver struct {
	a    *models.Account
	crud *db.CRUD
}

func (r *sysEditorResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

func (r *sysEditorResolver) Name() string {
	return r.a.Name
}

func (r *sysEditorResolver) Surname() string {
	return r.a.Surname
}

func (r *sysEditorResolver) Email() string {
	return r.a.Email
}

func (r *sysEditorResolver) RemoveRecruit(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()

	id := string(args.ID)

	// check that the ID is valid
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("id")
	}

	// attempt to remove question
	if err := r.crud.DeleteID(config.RecruitsCollection, bson.ObjectIdHex(id)); err != nil {
		return nil, er.NewGenericError()
	}
	result := "Recruit successfully removed."
	return &result, nil
}

func (r *sysEditorResolver) RemoveAccount(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()

	genericErr := "Failed to remove account."
	idStr := string(args.ID)
	if !bson.IsObjectIdHex(idStr) {
		return nil, er.NewInternalError(genericErr)
	}
	id := bson.ObjectIdHex(idStr)

	// check if there's an account with that id
	_, err := r.crud.FindOne(config.AccountsCollection, &bson.M{"_id": id})
	if err != nil {
		return nil, er.NewInternalError(genericErr)
	}

	// delete the account
	err = r.crud.DeleteID(config.AccountsCollection, id)
	if err != nil {
		log.Println("Failed to delete Account =>", err)
		return nil, er.NewGenericError()
	}

	// find the account's token manager
	rawTokenMgr, err := r.crud.FindOne(config.TokenManagersCollection, &bson.M{"account_id": id})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	// delete the account's token manager
	tokenMgr := transformTokenManager(rawTokenMgr)
	err = r.crud.DeleteID(config.TokenManagersCollection, tokenMgr.ID)
	if err != nil {
		log.Println("Failed to delete TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	msg := "Account successfully removed."
	return &msg, nil
}

// -----------------
// accountEditorResolver struct
// -----------------
type accountEditorResolver struct {
	a    *models.Account
	crud *db.CRUD
}

func (r *accountEditorResolver) RemoveAccount() (*string, error) {
	defer r.crud.CloseCopy()

	// delete the account
	err := r.crud.DeleteID(config.AccountsCollection, r.a.ID)
	if err != nil {
		log.Println("Failed to delete Account =>", err)
		return nil, er.NewGenericError()
	}

	// find the account's token manager
	rawTokenMgr, err := r.crud.FindOne(config.TokenManagersCollection, &bson.M{"account_id": r.a.ID})
	if err != nil {
		log.Println("Failed to find TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	// delete the account's token manager
	tokenMgr := transformTokenManager(rawTokenMgr)
	err = r.crud.DeleteID(config.TokenManagersCollection, tokenMgr.ID)
	if err != nil {
		log.Println("Failed to delete TokenManager =>", err)
		return nil, er.NewGenericError()
	}

	msg := "Account successfully removed."
	return &msg, nil
}
func (r *accountEditorResolver) CreateRecruit(args struct{ Info *recruitDetails }) (*recruitResolver, error) {
	// check if the account has a recruit profile
	account := r.a
	if !utils.IsNullID(account.RecruitID) {
		return nil, er.NewInputError("Account already has a Recruit profile.")
	}

	// check if info is nil
	info := args.Info
	if info == nil {
		return nil, er.NewMissingFieldError("info")
	}

	// validate info
	if info.Province == nil {
		return nil, er.NewMissingFieldError("info.province")
	}
	if info.Phone == nil {
		return nil, er.NewMissingFieldError("info.phone")
	}
	if info.Email == nil {
		return nil, er.NewMissingFieldError("info.email")
	}
	if info.City == nil {
		return nil, er.NewMissingFieldError("info.city")
	}
	if info.Gender == nil {
		return nil, er.NewMissingFieldError("info.gender")
	}
	if info.Disability == nil {
		return nil, er.NewMissingFieldError("info.disability")
	}
	if info.Vid1Url == nil {
		return nil, er.NewMissingFieldError("info.vid1_url")
	}
	if info.Vid2Url == nil {
		return nil, er.NewMissingFieldError("info.vid2_url")
	}
	if info.BirthYear == nil {
		return nil, er.NewMissingFieldError("info.birth_year")
	}
	if info.Qa1Question == nil {
		return nil, er.NewMissingFieldError("info.qa1_question")
	}
	if info.Qa1Answer == nil {
		return nil, er.NewMissingFieldError("info.qa1_answer")
	}
	if info.Qa2Question == nil {
		return nil, er.NewMissingFieldError("info.qa2_question")
	}
	if info.Qa2Answer == nil {
		return nil, er.NewMissingFieldError("info.qa2_answer")
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
		return nil, er.NewGenericError()
	}

	// attach the recruit profile to the account
	if err := r.crud.UpdateID(config.AccountsCollection, account.ID, bson.M{
		"recruit_id": recruit.ID,
	}); err != nil {
		log.Println(err)
		return nil, er.NewGenericError()
	}
	return &recruitResolver{&recruit, account}, nil
}

// -----------------
// editorResolver struct
// -----------------
type editorResolver struct {
	editor
}

func (r *editorResolver) ToRecruitEditor() (*recruitEditorResolver, bool) {
	v, ok := r.editor.(*recruitEditorResolver)
	return v, ok
}

func (r *editorResolver) ToSysEditor() (*sysEditorResolver, bool) {
	v, ok := r.editor.(*sysEditorResolver)
	return v, ok
}

func (r *editorResolver) ToAccountEditor() (*accountEditorResolver, bool) {
	v, ok := r.editor.(*accountEditorResolver)
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
