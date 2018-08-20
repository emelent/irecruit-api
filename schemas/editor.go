package schemas

// EditorSchema schema
var EditorSchema = Schema{
	Types: `
		union Editor = RecruitEditor | SysEditor | AccountEditor

		type RecruitEditor{
			removeRecruit: String
			#updateRecruit(): Recruit
		}
		
		type SysEditor{
			removeAccount(id: ID!): String
			
			#createQuestion(): Question
			#createIndustry(): Industry

			#removeRecruit(id: ID!): Recruit
			#removeAccount(id: ID!): Account
			#removeIndustry(id: ID!): Industry
			#removeQuestion(id: ID!): Industry

			#updateRecruit(id: ID!): Recruit
			#updateAccount(id: ID!): Account
			#updateIndustry(id: ID!): Industry
			#updateQuestion(id: ID!): Industry
		}

		type AccountEditor{
			createRecruit(info: RecruitDetails!): Recruit
			removeAccount(): String
			#updateAccount(): Account
		}
	`,
	Queries: `
	`,
	Mutations: `
		edit(token: String!, enforce: Enforce):Editor
	`,
}
