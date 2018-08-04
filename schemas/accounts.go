package schemas

// AccountSchema schema
var AccountSchema = Schema{
	Types: `
		type Fail {
			error: String!
		}

		type Account {
			id: ID!
			name: String!
			surname: String!
			email: String!
			#access_level: Int!
			hunter_id: ID
			recruit_id: ID
		}	
		
		type Tokens {
			refreshToken: String!
			accessToken: String!
		}

		union AccountOrFail = Account | Fail
		union TokensOrFail = Tokens | Fail
		
		input AccountDetails{
			email: String!
			password: String!
			name: String!
			surname: String!
		}


	`,
	Queries: `
		# Retrieve all accounts
		accounts(name: String): [Account]!
	`,
	Mutations: `
		removeAccount(id: ID!): AccountOrFail
		createAccount(info: AccountDetails!): TokensOrFail
	`,
}
