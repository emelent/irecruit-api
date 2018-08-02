package schemas

// LoginSchema sdl
var LoginSchema = Schema{
	Types:   ``,
	Queries: ``,
	Mutations: `
		login(email: String!, password: String!): String!
		#logout(token: String!): String!
	`,
}
