package routing

import (
	"net/http"

	db "../database"
	mware "../middleware"
	"github.com/gorilla/mux"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	resolver "../resolvers"
	schemas "../schemas"
)

// NewGqlHandler creates a graphql handler
func NewGqlHandler(crud *db.CRUD) http.Handler {
	// prepare graphql root resolver
	gqlResolver := &resolver.RootResolver{}
	gqlResolver.Init(crud)

	// create schema
	schema := graphql.MustParseSchema(
		schemas.CreateSchema(schemas.DefaultSchemas...),
		gqlResolver,
	)

	// make handler
	return mware.ReqInfoMiddleware(&relay.Handler{Schema: schema})
}

// NewGqlRouter prepares a new router with necessary endpoints
func NewGqlRouter(crud *db.CRUD, middleware ...mware.Middleware) http.Handler {
	// prepare router
	router := mux.NewRouter()
	mware.ApplyMiddleware(router, middleware...)

	// attach graphql handler
	router.
		Path("/graphql").
		// Methods(http.MethodPost). // use post-only in production
		Handler(NewGqlHandler(crud))

	return router
}
