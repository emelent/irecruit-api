package resolvers

import (
	"log"

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
		return nil, er.Generic()
	}

	// process results
	results := make([]*IndustryResolver, 0)
	for _, raw := range rawIndustries {
		industry := models.TransformIndustry(raw)
		results = append(results, &IndustryResolver{&industry})
	}
	return results, err
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
