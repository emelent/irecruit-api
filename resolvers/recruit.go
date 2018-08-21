package resolvers

import (
	"log"
	"time"

	config "../config"
	er "../errors"
	models "../models"
	utils "../utils"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Root Resolver methods
// -----------------

// Recruits resolves recruits gql method
func (r *RootResolver) Recruits() ([]*RecruitResolver, error) {
	defer r.crud.CloseCopy()

	results := make([]*RecruitResolver, 0)
	// get recruit profiles
	rawRecruits, err := r.crud.FindAll(config.RecruitsCollection, nil)
	if err != nil {
		log.Println(err)
		return results, er.Generic()
	}

	// process results
	for _, raw := range rawRecruits {
		var account models.Account
		recruit := models.TransformRecruit(raw)
		rawAccount, e := r.crud.FindOne(config.AccountsCollection, &bson.M{
			"recruit_id": recruit.ID,
		})
		if e == nil {
			account = models.TransformAccount(rawAccount)
			results = append(results, &RecruitResolver{&recruit, &account})
		}
	}
	return results, err
}

// CreateRecruit resolves createRecruit gql method
func (r *RootResolver) CreateRecruit(args struct {
	AccountID graphql.ID
	Info      *recruitDetails
}) (*RecruitResolver, error) {
	// check if id is valid
	id := string(args.AccountID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.InvalidField("id")
	}

	// check if there's an account with that id, i.e. retrieve the account dummy
	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(id))
	if err != nil {
		return nil, er.InvalidField("id")
	}

	// check if the account has a recruit profile
	account := models.TransformAccount(rawAccount)
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
	return &RecruitResolver{&recruit, &account}, nil
}

// RemoveRecruit resolves "removeRecruit" mutation
func (r *RootResolver) RemoveRecruit(args struct{ ID graphql.ID }) (*string, error) {
	return ResolveRemoveByID(
		r.crud,
		config.RecruitsCollection,
		"Recruit",
		string(args.ID),
	)
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
