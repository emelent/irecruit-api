package schemas

import (
	"fmt"
)

// Schema struct holds gql schema strings
type Schema struct {
	Types     string
	Queries   string
	Mutations string
}

// These are type definitions for interfaces and other type abstractions
var intermediateTypes = `
	type Fail {
		error: String!
	}
`

// DefaultSchemas is a list of all necessary schemas
var DefaultSchemas = []Schema{
	AccountSchema,
	LoginSchema,
	RecruitSchema,
	IndustrySchema,
}

// CreateSchema creates a schema from given Schema structs
func CreateSchema(schemas ...Schema) string {
	types := ""
	queries := ""
	mutations := ""
	for _, schema := range schemas {
		types += schema.Types + "\n"
		queries += schema.Queries + "\n"
		mutations += schema.Mutations + "\n"
	}
	schema := fmt.Sprintf(`
		schema {
			query: Query
			mutation: Mutation
		}

		type Query {
			%s
		}
		type Mutation{
			%s
		}

		%s

		%s
	`, queries, mutations, intermediateTypes, types)
	return schema
}
