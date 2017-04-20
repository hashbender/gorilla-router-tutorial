package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var (
	appContext AppContext
)

//Route this struct is used for declaring a route
type Route struct {
	Name             string
	Method           []string
	Pattern          string
	ContextedHandler *ContextedHandler
}

//Routes just stores our Route declarations
type Routes []Route

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

		}
	}
}

//NewRouter returns a new Gorrila Mux router
func NewRouter(c AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	appContext = c

	for _, route := range routes {
		//Check all routes to make sure the users are properly authenticated
		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(CheckAuth(SetContentTypeText(route.ContextedHandler)))
	}

	return router
}

//CheckAuth is an example middleware which demonstrates how we *might* check auth.
func CheckAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking Auth")
		//TODO this is just an example.
		//Get the cookie or something and check it
		cookie, err := r.Cookie("session")
		if err != nil || cookie == nil {
			// TODO if an err, then redirect
			// http.Redirect(w, r, "/", 401)
		}
		//If the auth check passes, then handle continue down the chain
		h.ServeHTTP(w, r)
	})
}

//SetContentTypeText this only exists to demonstrate how we can chain middlewares
func SetContentTypeText(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Setting Headers")
		w.Header().Set("Content-Type", "text/plain")
		h.ServeHTTP(w, r)
	})
}

//HelloWorldHandler just prints Hello World on a GET request
func HelloWorldHandler(c *AppContext, w http.ResponseWriter, req *http.Request) (int, error) {
	//So in this handler we now have the context provided
	fmt.Fprint(w, "Hello World")
	return http.StatusOK, nil
}

//GoodbyeWorldHandler just prints a Goodbye World on a GET request
func GoodbyeWorldHandler(c *AppContext, w http.ResponseWriter, req *http.Request) (int, error) {
	//So in this handler we now have the context provided
	fmt.Fprint(w, "Goodbye World")
	return http.StatusOK, nil
}

func main() {
	context := AppContext{Db: initDb(), Redis: initRedis()}

	http.ListenAndServe(":5000", NewRouter(context))
}

//TODO do your DB initialization here
func initDb() *sql.DB {
	return nil
}

//TODO do your redis pool initialization here
func initRedis() *redis.Pool {
	return nil
}

//Declare your routes and handlers here
var routes = Routes{
	Route{
		"HelloWorld",
		//You can handle more than just GET requests here
		[]string{"GET"},
		"/hello",
		&ContextedHandler{&appContext, HelloWorldHandler},
	},
	Route{
		"GoodbyeWorld",
		[]string{"GET"},
		"/goodbye",
		&ContextedHandler{&appContext, GoodbyeWorldHandler},
	},
}
