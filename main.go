package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	config "./config"
	db "./database"
	mware "./middleware"
	moc "./mocks"
	route "./routing"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	// seed the rand
	rand.Seed(time.Now().UnixNano())

	// get cli of args
	mockIt := flag.Bool("mock", false, "Uses mock database if true.")
	flag.Parse()

	// Setup environment
	config.SetupEnv()

	// setup crud system
	var crud *db.CRUD
	if *mockIt {
		// prepare mock
		log.Println("Using temporary mock db.")
		crud = moc.NewLoadedCRUD()
	} else {
		// create mongo url
		url := db.CreateMongoURL(
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		)
		log.Println("Dialing mongodb at =>", url)

		// connect to mongo
		mongoSession, err := mgo.Dial(url)
		if err != nil {
			log.Println("Mongodb didn't pick up, falling back to temporary mock database.")
		}
		crud = db.NewCRUD(mongoSession)
	}

	// closes db connection if any
	defer crud.Close()

	// handle unexpected panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovery in main()=>", r)
			os.Exit(1)
		}
	}()

	// prepare the router
	router := route.NewGqlRouter(crud, mware.CorsMiddleware, mware.LoggerMiddleware)

	// start listening
	port := ":" + string(os.Getenv("PORT"))
	log.Printf("Serving on 0.0.0.0 %s\n\n", port)
	log.Fatal(http.ListenAndServe(port, router))

}
