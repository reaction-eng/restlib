package routing

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/**
Create our own router that pulls in the mux router
*/
type Router struct {
	// Store the mux router
	*mux.Router

	//Store the paths so we can use them
	routes []Route
}

/**
Simple interface to allow other to specify routes
*/
type RouteProducer interface {
	/*
	 * Get any required routes
	 */
	GetRoutes() []Route
}

/**
* Build a new instance of this router.  It contains all of the paths so we can ghceck them later
 */
func NewRouter(optionsHandler http.HandlerFunc, routes []Route, routeProducers ...RouteProducer) *Router {
	muxRouter := mux.NewRouter().StrictSlash(true)

	//Add in an option to handle all options
	if optionsHandler != nil {
		muxRouter.Methods("OPTIONS").Handler(optionsHandler)
	}

	//Combine the newrouter into this one
	router := Router{
		muxRouter,
		make([]Route, 0),
	}

	//For each route
	for _, route := range routes {
		router.addRoute(route)

	}

	// Now march over each route produce
	for _, producer := range routeProducers {
		//For each route produced
		for _, route := range producer.GetRoutes() {
			router.addRoute(route)
		}

	}

	// Return pointer to router
	return &router
}

/**
Determines if it is a public path based upon the routes
*/
func (router *Router) addRoute(route Route) {

	// Define a new handler
	var handler http.Handler = route.HandlerFunc

	// Wrap the handler in a Logger wrapper
	handler = Logger(handler, route.Name)

	// Add the route to the router
	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)

	// Store this route so we can use it later
	router.routes = append(router.routes, route)

}

/**
Determines if it is a public path based upon the routes
*/
func (router *Router) PublicRoute(path string) bool {

	//Check each route
	for _, r := range router.routes {
		if r.Pattern == path {
			return r.Public
		}
	}

	//Assume not public by default
	return false
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		//See if there is a user
		userId := r.Context().Value("user") //Grab the id of the user that send the request

		//Make a user ID string
		userIdString := fmt.Sprint("userId", userId)

		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			userIdString,
			time.Since(start),
		)
	})
}