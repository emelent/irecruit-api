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

// DefaultSchemas is a list of all necessary schemas
var DefaultSchemas = []Schema{
	PublicSchema,
	ViewerSchema,
	DocumentSchema,
	EditorSchema,
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
	`, queries, mutations, types)
	return schema
}
