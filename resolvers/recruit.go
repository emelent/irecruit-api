package resolvers

import (
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
)

const recruitsCollection = "recruits"

type recruitDetails struct {
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
}

// QA Resolver
// type qaResolver struct {
// 	qa *models.QA
// }

// func (r *qaResolver) Question() string {
// 	return r.qa.Question
// }

// func (r *qaResolver) Answer() string {
// 	return r.qa.Answer
// }

// Recruit Resolver
type recruitResolver struct {
	r *models.Recruit
	a *models.Account
}

func (r *recruitResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
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

// func (r *recruitResolver) Qa1() *qaResolver {
// 	return &qaResolver{&r.r.Qa1}
// }

// func (r *recruitResolver) Qa2() *qaResolver {
// 	return &qaResolver{&r.r.Qa2}
// }

// Recruits resolves recruits gql method
func (r *RootResolver) Recruits() ([]*recruitResolver, error) {
	defer r.crud.CloseCopy()
	rawRecruits, err := r.crud.FindAll(recruitsCollection, nil)
	results := make([]*recruitResolver, 0)
	var account models.Account
	for _, r := range rawRecruits {
		recruit := transformRecruit(r)
		results = append(results, &recruitResolver{&recruit, &account})
	}
	return results, err
}

// CreateRecruit resolves createRecruit gql method
func (r *RootResolver) CreateRecruit(args struct {
	AccountID graphql.ID
	Info      *recruitDetails
}) (*recruitResolver, error) {
	return nil, nil
}
