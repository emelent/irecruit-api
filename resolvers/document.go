package resolvers

import (
	"log"

	config "../config"
	er "../errors"
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

type documentResolver struct {
	q *models.Document
}

func (r *documentResolver) ID() graphql.ID {
	return graphql.ID(r.q.ID.Hex())
}

func (r *documentResolver) OwnerID() graphql.ID {
	return graphql.ID(r.q.OwnerID.Hex())
}

func (r *documentResolver) OwnerType() string {
	return r.q.OwnerType
}

func (r *documentResolver) URL() string {
	return r.q.URL
}

func (r *documentResolver) DocType() string {
	return r.q.DocType
}

// Documents resolves "documents" gql query
func (r *RootResolver) Documents() ([]*documentResolver, error) {

	results := make([]*documentResolver, 0)
	// get documents
	rawDocuments, err := r.crud.FindAll(config.DocumentsCollection, nil)
	if err != nil {
		log.Println(err)
		return results, er.NewGenericError()
	}

	// process results
	for _, raw := range rawDocuments {
		document := transformDocument(raw)
		results = append(results, &documentResolver{&document})
	}
	return results, err
}

// CreateDocument resolves "createDocument"  gql mutation
func (r *RootResolver) CreateDocument(args struct {
	OwnerID   graphql.ID
	URL       string
	DocType   string
	OwnerType string
}) (*documentResolver, error) {
	defer r.crud.CloseCopy()

	// check that OwnerID is valid
	id := string(args.OwnerID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("owner_id")
	}

	// create document
	var document models.Document
	document.ID = bson.NewObjectId()
	document.URL = args.URL
	document.OwnerType = args.OwnerType
	document.DocType = args.DocType
	document.OwnerID = bson.ObjectIdHex(id)

	// validate document
	if err := document.OK(); err != nil {
		return nil, err
	}

	// attempt to insert
	if err := r.crud.Insert(config.DocumentsCollection, document); err != nil {
		return nil, er.NewGenericError()
	}

	return &documentResolver{&document}, nil
}

// RemoveDocument resolves "removeDocument" mutation
func (r *RootResolver) RemoveDocument(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()

	id := string(args.ID)

	// check that the ID is valid
	if !bson.IsObjectIdHex(id) {
		return nil, er.NewInvalidFieldError("id")
	}

	// attempt to remove document
	if err := r.crud.DeleteID(config.DocumentsCollection, bson.ObjectIdHex(id)); err != nil {
		return nil, er.NewGenericError()
	}
	result := "Document successfully removed."
	return &result, nil
}
