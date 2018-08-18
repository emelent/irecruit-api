package resolvers

import (
	"log"

	config "../config"
	er "../errors"
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

type questionResolver struct {
	q *models.Question
}

func (r *questionResolver) ID() graphql.ID {
	return graphql.ID(r.q.ID.Hex())
}

func (r *questionResolver) IndustryID() graphql.ID {
	return graphql.ID(r.q.IndustryID.Hex())
}

func (r *questionResolver) Question() string {
	return r.q.Question
}

// Questions resolves "question" gql query
func (r *RootResolver) Questions() ([]*questionResolver, error) {

	results := make([]*questionResolver, 0)
	// get industries
	rawQuestions, err := r.crud.FindAll(config.QuestionsCollection, nil)
	if err != nil {
		log.Println(err)
		return results, er.NewGenericError()
	}

	// process results
	for _, raw := range rawQuestions {
		question := transformQuestion(raw)
		results = append(results, &questionResolver{&question})
	}
	return results, err
}

// CreateQuestion resolves "createQuestion"  gql mutation
func (r *RootResolver) CreateQuestion(args struct {
	IndustryID graphql.ID
	Question   string
}) (*questionResolver, error) {
	defer r.crud.CloseCopy()

	// check that IndustryID is valid
	id := string(args.IndustryID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("industry_id")
	}

	var question models.Question
	question.ID = bson.NewObjectId()
	question.Question = args.Question
	question.IndustryID = bson.ObjectIdHex(id)

	// validate question
	if err := question.OK(); err != nil {
		return nil, err
	}

	// attempt to insert
	if err := r.crud.Insert(config.QuestionsCollection, question); err != nil {
		return nil, er.NewGenericError()
	}

	return &questionResolver{&question}, nil
}

// RemoveQuestion resolves "removeQuestion" mutation
func (r *RootResolver) RemoveQuestion(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()

	id := string(args.ID)

	// check that the ID is valid
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("id")
	}

	// attempt to remove question
	if err := r.crud.DeleteID(config.QuestionsCollection, bson.ObjectIdHex(id)); err != nil {
		return nil, er.NewGenericError()
	}
	result := "Question successfully removed."
	return &result, nil
}
