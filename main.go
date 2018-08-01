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

	if err := gqlResolver.OpenMongoDb(); err != nil {
		log.Fatal("Failed to open MongoDb connection.")
		return
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
		Handler(&relay.Handler{Schema: schema})

	port := ":9999"
	log.Printf("Serving on 0.0.0.0 %s\n\n", port)
	log.Fatal(http.ListenAndServe(port, router))

}
