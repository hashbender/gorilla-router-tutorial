package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

//ContextedHandler is a wrapper to provide AppContext to our Handlers
type ContextedHandler struct {
	*AppContext
	//ContextedHandlerFunc is the interface which our Handlers will implement
	ContextedHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
}

//AppContext provides the app context to handlers.  This *cannot* contain request-specific keys like
//sessionId or similar.  It is shared across requests.
type AppContext struct {
	Db    *sql.DB
	Redis *redis.Pool
}

func (handler ContextedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := handler.ContextedHandlerFunc(handler.AppContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		// TODO you can handle any error that a router might return here.
		// I use SendGrid to send an email anytime a 503 occurs :)
		}
	}
}

//HelloWorldHandler just prints Hello World on a GET request
func HelloWorldHandler(c *AppContext, w http.ResponseWriter, req *http.Request) (int, error) {
	//So in this handler we now have the context provided
	fmt.Fprint(w, "Hello World")
	return http.StatusOK, nil
}

//TODO do your DB initialization here
func initDb() *sql.DB {
	return nil
}

//TODO do your redis pool initialization here
func initRedis() *redis.Pool {
	return nil
}

func main() {
	context := AppContext{Db: initDb(), Redis: initRedis()}

	//Here we instantiate a ContextedHandler, which staisfies the ServeHTTP
	//interface and can thus be used as a Handler on a router instance
	contextedHandler := &ContextedHandler{&context, HelloWorldHandler}
	router := mux.NewRouter()
	router.Methods("GET").Path("/hello").Name("HelloWorldHandler").Handler(contextedHandler)
	http.ListenAndServe(":5000", router)
}
