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
		return nil, er.NewGenericError()
	}

	// process results
	results := make([]*QuestionResolver, 0)
	for _, raw := range rawQuestions {
		question := TransformQuestion(raw)
		results = append(results, &QuestionResolver{&question})
	}
	return results, err
}

// CreateQuestion resolves "createQuestion"  gql mutation
func (r *RootResolver) CreateQuestion(args struct {
	IndustryID graphql.ID
	Question   string
}) (*QuestionResolver, error) {
	defer r.crud.CloseCopy()

	// check that IndustryID is valid
	id := string(args.IndustryID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("industry_id")
	}

	// create question
	var question models.Question
	question.ID = bson.NewObjectId()
	question.Question = args.Question
	question.IndustryID = bson.ObjectIdHex(id)

	// validate question
	if err := question.OK(); err != nil {
		return nil, err
	}

	// store question in db
	if err := r.crud.Insert(config.QuestionsCollection, question); err != nil {
		return nil, er.NewGenericError()
	}

	// return question
	return &QuestionResolver{&question}, nil
}

// RemoveQuestion resolves "removeQuestion" mutation
func (r *RootResolver) RemoveQuestion(args struct{ ID graphql.ID }) (*string, error) {
	return ResolveRemoveByID(
		r.crud,
		config.QuestionsCollection,
		"Question",
		string(args.ID),
	)
}

// RandomQuestions resolves "randomQuestions" gql query
func (r *RootResolver) RandomQuestions(args struct{ IndustryID graphql.ID }) ([]*QuestionResolver, error) {
	defer r.crud.CloseCopy()

	// check that the ID is valid
	id := string(args.IndustryID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("id")
	}

	// get industries
	rawQuestions, err := r.crud.FindAll(config.QuestionsCollection, &bson.M{"industry_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println(err)
		return nil, er.NewGenericError()
	}

	// process results
	randomQuestions := make([]*QuestionResolver, 0)
	rawQuestions = utils.PickRandomN(2, rawQuestions)
	for _, raw := range rawQuestions {
		question := TransformQuestion(raw)
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
