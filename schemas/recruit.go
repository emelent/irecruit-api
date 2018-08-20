package schemas

// RecruitSchema graphql schema
var RecruitSchema = Schema{
	Types: `
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
			qa1_question: String
			qa1_answer: String
			qa2_question: String
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
		recruits: [Recruit]!
	`,
	Mutations: `
		createRecruit(account_id: ID!, info: RecruitDetails!): Recruit
		removeRecruit(id: ID!): String

	`,
}
