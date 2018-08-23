package schemas

// IndustrySchema schema
var IndustrySchema = Schema{
	Types: `
		type Industry{
			id: ID!
			name: String!
		}
	`,
	Queries: `
		industries:[Industry]!
	`,
	Mutations: `
	`,
}
