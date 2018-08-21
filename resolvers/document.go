package resolvers

import (
	"log"

	config "../config"
	er "../errors"
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

// -----------------
// Root Resolver methods
// -----------------

// Documents resolves "documents" gql query
func (r *RootResolver) Documents() ([]*DocumentResolver, error) {
	defer r.crud.CloseCopy()
	// get documents
	rawDocuments, err := r.crud.FindAll(config.DocumentsCollection, nil)
	if err != nil {
		log.Println(err)
		return nil, er.Generic()
	}

	// process results
	results := make([]*DocumentResolver, 0)
	for _, raw := range rawDocuments {
		document := models.TransformDocument(raw)
		results = append(results, &DocumentResolver{&document})
	}
	return results, err
}

// CreateDocument resolves "createDocument"  gql mutation
func (r *RootResolver) CreateDocument(args struct {
	OwnerID   graphql.ID
	URL       string
	DocType   string
	OwnerType string
}) (*DocumentResolver, error) {
	defer r.crud.CloseCopy()

	// check that OwnerID is valid
	id := string(args.OwnerID)
	if !bson.IsObjectIdHex(id) {
		return nil, er.InvalidField("owner_id")
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
		return nil, er.Generic()
	}

	return &DocumentResolver{&document}, nil
}

// RemoveDocument resolves "removeDocument" mutation
func (r *RootResolver) RemoveDocument(args struct{ ID graphql.ID }) (*string, error) {
	return ResolveRemoveByID(
		r.crud,
		config.DocumentsCollection,
		"Document",
		string(args.ID),
	)
}

// -----------------
// DocumentResolver struct
// -----------------

// DocumentResolver resolves Document
type DocumentResolver struct {
	q *models.Document
}

// ID resolves Document.ID
func (r *DocumentResolver) ID() graphql.ID {
	return graphql.ID(r.q.ID.Hex())
}

// OwnerID resolves Document.OwnerID
func (r *DocumentResolver) OwnerID() graphql.ID {
	return graphql.ID(r.q.OwnerID.Hex())
}

// OwnerType resolves Document.OwnerType
func (r *DocumentResolver) OwnerType() string {
	return r.q.OwnerType
}

// URL resolves Document.URL
func (r *DocumentResolver) URL() string {
	return r.q.URL
}

// DocType resolves Document.DocType
func (r *DocumentResolver) DocType() string {
	return r.q.DocType
}
