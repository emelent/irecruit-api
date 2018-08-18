package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/mgo.v2/bson"

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

// NewUploadHandler creates an upload handler
func NewUploadHandler() func(http.ResponseWriter, *http.Request) {
	const maxUploadSize = 100 * 1024
	fmt.Println(maxUploadSize)
	const uploadPath = "./files"
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			log.Println(err)
			jsonEncode(w, "File too large.", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			log.Println(err)
			jsonEncode(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			jsonEncode(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		filetype := http.DetectContentType(fileBytes)
		if filetype != "image/jpeg" && filetype != "image/jpg" &&
			filetype != "image/gif" && filetype != "image/png" &&
			filetype != "application/pdf" {
			jsonEncode(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}

		fileName := bson.NewObjectId().Hex()
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			log.Println(err)
			jsonEncode(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		fmt.Println(fileEndings)
		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", fileEndings[0], newPath)

		newFile, err := os.Create(newPath)
		if err != nil {
			log.Println(err)
			jsonEncode(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			log.Println(err)
			jsonEncode(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	}
}

// jsonEncode writes a json response
func jsonEncode(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewRouter prepares a new router with necessary endpoints
func NewRouter(crud *db.CRUD, middleware ...mware.Middleware) http.Handler {
	// prepare router
	router := mux.NewRouter()
	mware.ApplyMiddleware(router, middleware...)

	// attach graphql handler
	router.
		Path("/graphql").
		// Methods(http.MethodPost). // strictly post for production in production
		Handler(NewGqlHandler(crud))

	// attach upload handler
	router.
		Path("/upload").
		Methods(http.MethodPost).
		HandlerFunc(NewUploadHandler())

	// attach static file  handler
	dir := http.Dir("./files")
	router.
		PathPrefix("/file/").
		Handler(http.StripPrefix("/file/", http.FileServer(dir))).
		Methods(http.MethodGet)

	return router
}
