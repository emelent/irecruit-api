package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

//Middleware func
type Middleware func(http.Handler) http.Handler

//LoggerMiddleware logs each request
func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

//CorsMiddleware allows cross origin access
func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		h.ServeHTTP(w, r)
	})
}

// CK get
type CK string

//UaKey is lsf
const UaKey = CK("user_agent")

//IPKey is asdf
const IPKey = CK("ip_address")

// ReqInfoMiddleware puts some request info into the context
func ReqInfoMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), UaKey, r.Header.Get("User-Agent"))
		ctx = context.WithValue(ctx, IPKey, r.RemoteAddr)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

//ApplyMiddleware   applies given middleware to router
func ApplyMiddleware(router http.Handler, middleware ...Middleware) http.Handler {
	newRouter := router
	for _, m := range middleware {
		newRouter = m(newRouter)
	}
	return newRouter
}
