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

// QA Resolver
type qaResolver struct {
	qa *models.QA
}

func (r *qaResolver) Question() string {
	return r.qa.Question
}

func (r *qaResolver) Answer() string {
	return r.qa.Answer
}

// Recruit Resolver
type recruitResolver struct {
	r *models.Recruit
	a *models.Account
}

func (r *recruitResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
}

func (r *recruitResolver) Age() int32 {
	year := int32(time.Now().Year())
	return year - r.r.BirthYear
}

func (r *recruitResolver) Name() string {
	return r.a.Name
}

func (r *recruitResolver) Surname() string {
	return r.a.Surname
}

func (r *recruitResolver) Phone() string {
	return r.r.Phone
}

func (r *recruitResolver) Email() string {
	return r.r.Email
}

func (r *recruitResolver) Province() string {
	return r.r.Province
}

func (r *recruitResolver) City() string {
	return r.r.City
}

func (r *recruitResolver) Gender() string {
	return r.r.Gender
}

func (r *recruitResolver) Disability() string {
	return r.r.Disability
}

func (r *recruitResolver) Vid1Url() string {
	return r.r.Vid1Url
}

func (r *recruitResolver) Vid2Url() string {
	return r.r.Vid2Url
}

func (r *recruitResolver) Qa1() *qaResolver {
	return &qaResolver{&r.r.Qa1}
}

func (r *recruitResolver) Qa2() *qaResolver {
	return &qaResolver{&r.r.Qa2}
}

// Recruits resolves recruits gql method
func (r *RootResolver) Recruits() ([]*recruitResolver, error) {
	defer r.crud.CloseCopy()

	results := make([]*recruitResolver, 0)
	// get recruit profiles
	rawRecruits, err := r.crud.FindAll(config.RecruitsCollection, nil)
	if err != nil {
		log.Println(err)
		return results, er.NewGenericError()
	}

	// process results
	for _, raw := range rawRecruits {
		var account models.Account
		recruit := transformRecruit(raw)
		rawAccount, e := r.crud.FindOne(config.AccountsCollection, &bson.M{
			"recruit_id": recruit.ID,
		})
		if e == nil {
			account = transformAccount(rawAccount)
			results = append(results, &recruitResolver{&recruit, &account})
		}
	}
	return results, err
}

// CreateRecruit resolves createRecruit gql method
func (r *RootResolver) CreateRecruit(args struct {
	AccountID graphql.ID
	Info      *recruitDetails
}) (*recruitResolver, error) {
	// check if id is valid
	id := string(args.AccountID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("id")
	}

	// check if there's an account with that id, i.e. retrieve the account dummy
	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(id))
	if err != nil {
		return nil, er.NewInvalidFieldError("id")
	}

	// check if the account has a recruit profile
	account := transformAccount(rawAccount)
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
	return &recruitResolver{&recruit, &account}, nil
}

// RemoveRecruit resolves "removeRecruit" mutation
func (r *RootResolver) RemoveRecruit(args struct{ ID graphql.ID }) (*string, error) {
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
