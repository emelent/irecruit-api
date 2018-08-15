package resolvers

import (
	"log"

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
