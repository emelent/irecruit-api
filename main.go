package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	mware "./middleware"
	resolver "./resolvers"
	schemas "./schemas"
	"github.com/gorilla/mux"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	gqlResolver := &resolver.RootResolver{}

	gqlResolver.OpenMongoDb()
	if err := gqlResolver.OpenMongoDb(); err != nil {
		log.Println("MongoDb failed to connect, falling back to temporary mock database.")
	}
	defer gqlResolver.CloseMongoDb()

	schema := graphql.MustParseSchema(
		schemas.CreateSchema(schemas.DefaultSchemas...),
		gqlResolver,
	)
	router := mux.NewRouter()
	mware.ApplyMiddleware(router, mware.CorsMiddleware, mware.LoggerMiddleware)
	router.
		Path("/graphql").
		// Methods(http.MethodPost).
		Handler(mware.ReqInfoMiddleware(&relay.Handler{Schema: schema}))

	port := ":9999"
	log.Printf("Serving on 0.0.0.0 %s\n\n", port)
	log.Fatal(http.ListenAndServe(port, router))

}
