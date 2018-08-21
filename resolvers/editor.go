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

	editAsRecruit := func() (*EditorResolver, error) {
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

		// return RecruitEditor
		recruit := models.TransformRecruit(rawRecruit)
		Editor := &RecruitEditorResolver{&recruit, &account, r.crud}
		return &EditorResolver{Editor}, nil
	}
	editAsAccount := func() (*EditorResolver, error) {
		return &EditorResolver{&AccountEditorResolver{&account, r.crud}}, nil
	}

	editAsSys := func() (*EditorResolver, error) {
		// check if account is sys account
		if !utils.IsSysAccount(&account) {
			return nil, er.Input("Failed to enfore 'SYSTEM'.")
		}

		// return sysEditor
		Editor := &SysEditorResolver{&account, r.crud}
		return &EditorResolver{Editor}, nil
	}

	// enforce an enforceable
	if args.Enforce != nil {
		switch *args.Enforce {
		case "RECRUIT":
			return editAsRecruit()
		case "HUNTER": // try to enforce hunter
			return nil, er.Input("Unimplemented")

		case "SYSTEM": // try to enforce system
			return editAsSys()
		case "ACCOUNT":
			return editAsAccount()

		}
	}

	// try to edit as SysEditor
	if editor, err := editAsSys(); err == nil {
		return editor, nil
	}

	// try to edit as RecruitEditor
	if editor, err := editAsRecruit(); err == nil {
		return editor, nil
	}

	// if all else fails edit as AccountEditor
	return editAsAccount()
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

// UpdateRecruit resolves RecruitEditor.UpdateRecruit
func (r *RecruitEditorResolver) UpdateRecruit(args struct {
	Info *recruitDetails
}) (*RecruitResolver, error) {
	defer r.crud.CloseCopy()

	// prepare updates
	updates := bson.M{}
	info := args.Info
	if info.Phone != nil {
		updates["phone"] = *info.Phone
	}
	if info.Email != nil {
		updates["email"] = *info.Email
	}
	if info.Province != nil {
		updates["province"] = *info.Province
	}
	if info.City != nil {
		updates["city"] = *info.City
	}
	if info.Gender != nil {
		updates["gender"] = *info.Gender
	}
	if info.Disability != nil {
		updates["disability"] = *info.Disability
	}
	if info.Vid1Url != nil {
		updates["vid1_url"] = *info.Vid1Url
	}
	if info.Vid2Url != nil {
		updates["vid2_url"] = *info.Vid2Url
	}
	if info.BirthYear != nil {
		updates["birth_year"] = *info.BirthYear
	}

	// perform update
	rawRecruit, err := GenericUpdateByID(r.crud, config.RecruitsCollection, r.r.ID, updates)
	if err != nil {
		return nil, err
	}

	// return updated recruit profile
	recruit := models.TransformRecruit(rawRecruit)
	return &RecruitResolver{&recruit, r.a}, nil
}

// UpdateQAs resolves RecruitEditor.UpdateQAs
func (r *RecruitEditorResolver) UpdateQAs(args struct {
	Qa1 *qaDetails
	Qa2 *qaDetails
}) ([]*QaResolver, error) {
	defer r.crud.CloseCopy()

	// check that we have at least 1 QA
	qa1 := args.Qa1
	qa2 := args.Qa2
	if qa1 == nil && qa2 == nil {
		return nil, er.Input("No QAs given.")
	}

	updates := bson.M{}
	getQuestion := func(id bson.ObjectId) (*models.Question, error) {
		rawQ, err := r.crud.FindID(config.QuestionsCollection, id)
		if err != nil {
			return nil, err
		}
		question := models.TransformQuestion(rawQ)
		return &question, nil
	}

	results := make([]*QaResolver, 0)
	// prepare qa1 update
	if qa1 != nil {
		id := string(qa1.QuestionID)
		if !bson.IsObjectIdHex(id) {
			return nil, er.InvalidField("qa1.question_id")
		}
		question, err := getQuestion(bson.ObjectIdHex(id))
		if err != nil {
			return nil, er.InvalidField("qa1.question_id")
		}
		qa := models.QA{
			Question: question.Question,
			Answer:   qa1.Answer,
		}
		updates["qa1"] = qa
		results = append(results, &QaResolver{&qa})

	}

	// prepare qa2 update
	if qa2 != nil {
		id := string(qa2.QuestionID)
		if !bson.IsObjectIdHex(id) {
			return nil, er.InvalidField("qa2.question_id")
		}
		question, err := getQuestion(bson.ObjectIdHex(id))
		if err != nil {
			return nil, er.InvalidField("qa2.question_id")
		}
		qa := models.QA{
			Question: question.Question,
			Answer:   qa2.Answer,
		}
		updates["qa2"] = qa
		results = append(results, &QaResolver{&qa})
	}

	if _, err := GenericUpdateByID(
		r.crud,
		config.RecruitsCollection,
		r.r.ID,
		updates,
	); err != nil {
		return nil, er.Generic()
	}

	return results, nil
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

// UpdateAccount resolves AccountEditor.UpdateAccount
func (r *AccountEditorResolver) UpdateAccount(args struct {
	Info *accountDetails
}) (*AccountResolver, error) {
	defer r.crud.CloseCopy()

	// prepare updates
	updates := bson.M{}
	info := args.Info

	if info.Email != nil {
		updates["email"] = *info.Email
	}
	if info.Name != nil {
		updates["name"] = *info.Name
	}
	if info.Surname != nil {
		updates["surname"] = *info.Surname
	}

	// perform update
	rawAccount, err := GenericUpdateByID(r.crud, config.AccountsCollection, r.a.ID, updates)
	if err != nil {
		return nil, err
	}

	// return updated account
	account := models.TransformAccount(rawAccount)
	return &AccountResolver{&account}, nil
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

// -----------------
// qaDetails struct
// -----------------
type qaDetails struct {
	QuestionID graphql.ID
	Answer     string
}
