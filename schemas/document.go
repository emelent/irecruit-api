package schemas

// DocumentSchema graphql schema for documents
var DocumentSchema = Schema{
	Types: `
		type Document{
			id: ID!
			url: String!
			doc_type: DocType!
			owner_type: OwnerType!
			owner_id: ID!
		}

		enum  DocType{
			QUALIFICATION
			COMPANY
		}

		enum OwnerType{
			RECRUIT
			COMPANY
		}
	`,
	Queries: `
	`,
	Mutations: `
		createDocument(
			url: String!,
			doc_type: DocType!,
			owner_type: OwnerType!,
			owner_id:ID!
		): Document
	`,
}
