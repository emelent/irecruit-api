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
		}

	#	type Guest implements Viewer{
	#		id: ID!
	#		name: String!
	#		surname: String!
	#		email: String!
	#		industries: [Industry]!
	#	}
	`,
	Queries: `
		view(token: String):Viewer
	`,
	Mutations: `
	`,
}
