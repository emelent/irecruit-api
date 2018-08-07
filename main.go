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
	route "./routing"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	// seed the rand
	rand.Seed(time.Now().UnixNano())

	// get cli of args
	mock := flag.Bool("mock", false, "Uses mock database if true.")
	flag.Parse()

	// setup crud system
	var crud *db.CRUD
	if *mock != true {
		// connect to mongo
		log.Println("Dialing mongodb ...")
		mongoSession, err := mgo.Dial(config.DbHost)
		crud = db.NewCRUD(mongoSession)
		if err != nil {
			log.Println("Mongodb didn't pick up, falling back to temporary mock database.")
		}
	} else {
		log.Println("Using temporary mock db.")
		crud = db.NewCRUD(nil)
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
	log.Printf("Serving on 0.0.0.0 %s\n\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))

}
