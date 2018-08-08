package schemas

// LoginSchema sdl
var LoginSchema = Schema{
	Types:   ``,
	Queries: ``,
	Mutations: `
		login(email: String!, password: String!): Tokens
		#logout(token: String!): String!
	`,
}
