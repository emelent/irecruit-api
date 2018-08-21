package schemas

// ViewerSchema schema
var ViewerSchema = Schema{
	Types: `

		interface Viewer{
			id: ID!
			name: String!
			surname: String!
			email: String!
		}

		type AccountViewer implements Viewer{
			id: ID!
			name: String!
			surname: String!
			email:  String!
			is_hunter: Boolean!
			is_recruit:  Boolean!
			checkPassword(password: String!): Boolean!

		}

		type RecruitViewer implements Viewer{
			id: ID!
			name: String!
			surname: String!
			email: String!
			profile: Recruit
		}
		
		#type HunterViewer implements Viewer{
		#	id: ID!
		#	name: String!
		#	surname: String!
		#	email: String!
		#	recruit(id:ID!): Recruit
		#}

		type SysViewer implements Viewer{
			id: ID!
			name: String!
			surname: String!
			email: String!
			accounts: [Account]!
			#recruits: [Recruit]!
			#questions: [Question]!
			#documents: [Documents]!
		}

		enum Enforce{
			RECRUIT,
			HUNTER,
			SYSTEM,
			ACCOUNT
		}
	`,
	Queries: `
		view(token: String!, enforce: Enforce):Viewer
	`,
	Mutations: `
	`,
}
