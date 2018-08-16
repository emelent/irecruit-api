package resolvers

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	config "../config"
	er "../errors"
	models "../models"
	graphql "github.com/graph-gophers/graphql-go"
)

type industryResolver struct {
	i *models.Industry
}

func (r *industryResolver) ID() graphql.ID {
	return graphql.ID(r.i.ID.Hex())
}

func (r *industryResolver) Name() string {
	return r.i.Name
}

// Industries resolves "industries" gql query
func (r *RootResolver) Industries() ([]*industryResolver, error) {
	defer r.crud.CloseCopy()

	results := make([]*industryResolver, 0)
	// get industries
	rawIndustries, err := r.crud.FindAll(config.IndustriesCollection, nil)
	if err != nil {
		log.Println(err)
		return results, er.NewGenericError()
	}

	// process results
	for _, raw := range rawIndustries {
		industry := transformIndustry(raw)
		results = append(results, &industryResolver{&industry})
	}
	return results, err
}

// CreateIndustry resolves "createIndustry"  gql mutation
func (r *RootResolver) CreateIndustry(args struct{ Name string }) (*industryResolver, error) {
	defer r.crud.CloseCopy()

	var industry models.Industry
	industry.ID = bson.NewObjectId()
	industry.Name = args.Name

	// check that the name does not already exist
	if _, err := r.crud.FindOne(config.IndustriesCollection, &bson.M{
		"name": args.Name,
	}); err == nil {
		return nil, er.NewInputError("An industry by that name already exists.")
	}

	// validate industry
	if err := industry.OK(); err != nil {
		return nil, err
	}

	// attempt to insert
	if err := r.crud.Insert(config.IndustriesCollection, industry); err != nil {
		return nil, er.NewGenericError()
	}

	return &industryResolver{&industry}, nil
}

// RemoveIndustry resolves "removeIndustry" mutation
func (r *RootResolver) RemoveIndustry(args struct{ ID graphql.ID }) (*string, error) {
	defer r.crud.CloseCopy()
	return nil, nil
}
