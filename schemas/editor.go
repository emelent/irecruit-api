package schemas

// EditorSchema schema
var EditorSchema = Schema{
	Types: `
		union Editor = RecruitEditor | SysEditor | AccountEditor

		input QaDetails{
			question_id: ID!
			answer: String!
		}

		type AccountEditor{
			createRecruit(info: RecruitDetails!): Recruit
			removeAccount(): String
			updateAccount(info: AccountDetails): Account
		}

		type RecruitEditor{
			removeRecruit: String
			updateRecruit(info: RecruitDetails): Recruit
			updateQAs(qa1: QaDetails, qa2: QaDetails): [QA]!
		}
		
		type SysEditor{
			createQuestion(industry_id: ID!, question: String!): Question
			createIndustry(name: String!): Industry

			removeAccount(id: ID!): String
			removeRecruit(id: ID!): String
			removeIndustry(id: ID!): String
			removeQuestion(id: ID!): String
			removeDocument(id: ID!): String
			

			#updateRecruit(id: ID!): Recruit
			#updateAccount(id: ID!): Account
			#updateIndustry(id: ID!): Industry
			#updateQuestion(id: ID!): Industry
		}
	`,
	Queries: `
	`,
	Mutations: `
		edit(token: String!, enforce: Enforce):Editor
	`,
}
