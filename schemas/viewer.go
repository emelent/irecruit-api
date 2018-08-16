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
			profile: Recruit
		}
		
		#type HunterViewer implements Viewer{
		#	recruit(id:ID!): Recruit
		#}

		type SysViewer implements Viewer{
			accounts: [Account]!
			#recruits: [Recruit]!
		}

	#	type Guest implements Viewer{
	#		industries: [Industry]!
	#	}
	`,
	Queries: `
		view(token: String):Viewer
	`,
	Mutations: `
	`,
}
