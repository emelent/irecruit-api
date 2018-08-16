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
		createIndustry(name:String!):Industry
		#removeIndustry(id:ID!)
	`,
}
