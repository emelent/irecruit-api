package resolvers

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	config "../config"
	er "../errors"
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
)

// -----------------
// Root Resolver methods
// -----------------

// Industries resolves "industries" gql query
func (r *RootResolver) Industries() ([]*IndustryResolver, error) {
	defer r.crud.CloseCopy()

	// get industries
	rawIndustries, err := r.crud.FindAll(config.IndustriesCollection, nil)
	if err != nil {
		log.Println(err)
		return nil, er.NewGenericError()
	}

	// process results
	results := make([]*IndustryResolver, 0)
	for _, raw := range rawIndustries {
		industry := transformIndustry(raw)
		results = append(results, &IndustryResolver{&industry})
	}
	return results, err
}

// CreateIndustry resolves "createIndustry"  gql mutation
func (r *RootResolver) CreateIndustry(args struct{ Name string }) (*IndustryResolver, error) {
	defer r.crud.CloseCopy()

	// check that the name does not already exist
	if _, err := r.crud.FindOne(config.IndustriesCollection, &bson.M{
		"name": args.Name,
	}); err == nil {
		return nil, er.NewInputError("An industry by that name already exists.")
	}

	// create industry
	var industry models.Industry
	industry.ID = bson.NewObjectId()
	industry.Name = args.Name

	// validate industry
	if err := industry.OK(); err != nil {
		return nil, err
	}

	// store industry in db
	if err := r.crud.Insert(config.IndustriesCollection, industry); err != nil {
		return nil, er.NewGenericError()
	}

	// return industry
	return &IndustryResolver{&industry}, nil
}

// RemoveIndustry resolves "removeIndustry" mutation
func (r *RootResolver) RemoveIndustry(args struct{ ID graphql.ID }) (*string, error) {
	return ResolveRemoveByID(
		r.crud,
		config.IndustriesCollection,
		"Industry",
		string(args.ID),
	)
}

// -----------------
// IndustryResolver struct
// -----------------

// IndustryResolver resolves Industry
type IndustryResolver struct {
	i *models.Industry
}

// ID resolves Industry.ID
func (r *IndustryResolver) ID() graphql.ID {
	return graphql.ID(r.i.ID.Hex())
}

// Name resolves Industry.Name
func (r *IndustryResolver) Name() string {
	return r.i.Name
}
