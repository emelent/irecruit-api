package schemas

// PublicSchema public graphql queries
var PublicSchema = Schema{
	Types: `
		type Account {
			id: ID!
			name: String!
			surname: String!
			email: String!
			hunter_id: ID!
			recruit_id: ID!
		}	
		
		type Tokens {
			refreshToken: String!
			accessToken: String!
		}
		
		input AccountDetails{
			email: String
			password: String
			name: String
			surname: String
		}

		type Industry{
			id: ID!
			name: String!
		}
	`,
	Queries: `
		industries:[Industry]!
	`,
	Mutations: `
		createAccount(info: AccountDetails!): Tokens
		login(email: String!, password: String!): Tokens
		#logout(token: String!): String!
	`,
}
