package resolvers

import (
	"log"

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

// Questions resolves "questions" gql query
func (r *RootResolver) Questions() ([]*QuestionResolver, error) {
	defer r.crud.CloseCopy()

	// get industries
	rawQuestions, err := r.crud.FindAll(config.QuestionsCollection, nil)
	if err != nil {
		log.Println(err)
		return nil, er.Generic()
	}

	// process results
	results := make([]*QuestionResolver, 0)
	for _, raw := range rawQuestions {
		question := models.TransformQuestion(raw)
		results = append(results, &QuestionResolver{&question})
	}
	return results, err
}

// RandomQuestions resolves "randomQuestions" gql query
func (r *RootResolver) RandomQuestions(args struct{ IndustryID graphql.ID }) ([]*QuestionResolver, error) {
	defer r.crud.CloseCopy()

	// check that the ID is valid
	id := string(args.IndustryID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.InvalidField("id")
	}

	// get industries
	rawQuestions, err := r.crud.FindAll(config.QuestionsCollection, bson.M{"industry_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println(err)
		return nil, er.Generic()
	}

	// process results
	randomQuestions := make([]*QuestionResolver, 0)
	rawQuestions = utils.PickRandomN(2, rawQuestions)
	for _, raw := range rawQuestions {
		question := models.TransformQuestion(raw)
		randomQuestions = append(randomQuestions, &QuestionResolver{&question})
	}

	// return randomQuestions
	return randomQuestions, err
}

// -----------------
// QuestionResolver struct
// -----------------

// QuestionResolver resolves Question
type QuestionResolver struct {
	q *models.Question
}

// ID resolves Question.ID
func (r *QuestionResolver) ID() graphql.ID {
	return graphql.ID(r.q.ID.Hex())
}

// IndustryID resolves Question.IndustryID
func (r *QuestionResolver) IndustryID() graphql.ID {
	return graphql.ID(r.q.IndustryID.Hex())
}

// Question resolves Question.Question
func (r *QuestionResolver) Question() string {
	return r.q.Question
}
