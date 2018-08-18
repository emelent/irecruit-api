package schemas

// QuestionSchema graphql schema for Question
var QuestionSchema = Schema{
	Types: `
		type Question{
			id: ID!
			industry_id: ID!
			question: String!
		}
	`,
	Queries: `
		questions:[Question]!
		randomQuestions(industry_id: ID!): [Question]!
	`,
	Mutations: `
		createQuestion(industry_id: ID!, question: String!): Question
		removeQuestion(id: ID!): String
	`,
}
