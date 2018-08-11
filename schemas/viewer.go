package schemas

// ViewerSchema schema
var ViewerSchema = Schema{
	Types: `
		type Viewer = Hunter | Recruit | SysUser | Guest

		viewer(token: String): Viewer

		
	`,
	Queries: `

	`,
	Mutations: `
	`,
}
