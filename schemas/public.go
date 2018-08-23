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

		type Question{
			id: ID!
			industry_id: ID!
			question: String!
		}

		type QA{
			question: String!
			answer: String!
		}

		scalar Date
		
		type Recruit{
			id: ID!
			surname: String!
			phone: String!
			name: String!
			email: String!
			age: Int!
			province: String!
			city: String!
			gender: String!
			disability: String!
			vid1_url: String!
			vid2_url: String!		
			qa1: QA!
			qa2: QA!
		}

		input RecruitDetails{
			province: Province
			city: String
			gender: Gender
			disability: String
			vid1_url: String
			vid2_url: String
			phone: String
			email: String
			birth_year: Int
			qa1_question_id: ID
			qa1_answer: String
			qa2_question_id: ID
			qa2_answer: String
		}

		enum Province{
			KWAZULU_NATAL
			NORTHERN_CAPE
			WESTERN_CAPE
			EASTERN_CAPE
			NORTH_WEST
			FREE_STATE
			MPUMALANGA
			GAUTENG
			LIMPOPO
		}

		enum Gender{
			MALE
			FEMALE
		}
	`,
	Queries: `
		industries:[Industry]!
		randomQuestions(industry_id: ID!): [Question]!
	`,
	Mutations: `
		createAccount(info: AccountDetails!): Tokens
		login(email: String!, password: String!): Tokens
		#logout(token: String!): String!
	`,
}
