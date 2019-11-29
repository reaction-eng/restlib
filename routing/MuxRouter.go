// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

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
type MuxRouter struct {
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

//GetType def a logger wrapper
type LoggerWrapper func(inner http.Handler, name string) http.Handler

/**
* Build a new instance of this router.  It contains all of the paths so we can ghceck them later
 */
func NewRouter(optionsHandler http.HandlerFunc, routes []Route, loggerWrapper LoggerWrapper, routeProducers ...RouteProducer) *MuxRouter {
	muxRouter := mux.NewRouter().StrictSlash(true)

	//Add in an option to handle all options
	if optionsHandler != nil {
		muxRouter.Methods("OPTIONS").Handler(optionsHandler)
	}

	//Combine the newrouter into this one
	router := MuxRouter{
		muxRouter,
		make([]Route, 0),
	}

	//For each route
	for _, route := range routes {
		router.addRoute(route, loggerWrapper)

	}

	// Now march over each route produce
	for _, producer := range routeProducers {
		//For each route produced
		for _, route := range producer.GetRoutes() {
			router.addRoute(route, loggerWrapper)
		}

	}

	// Return pointer to router
	return &router
}

/**
Determines if it is a public path based upon the routes
*/
func (router *MuxRouter) addRoute(route Route, loggerWrapper LoggerWrapper) {

	// Define a new handler
	var handler http.Handler = route.HandlerFunc

	// Wrap the handler in a Logger wrapper
	if loggerWrapper != nil {
		handler = loggerWrapper(handler, route.Name)
	}
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
func (router *MuxRouter) GetRoute(req *http.Request) *Route {

	//Get the route name from the req
	if muxRoute := mux.CurrentRoute(req); muxRoute != nil {
		routeName := muxRoute.GetName()

		//Now march over the routes and return one
		for _, route := range router.routes {
			if route.Name == routeName {
				return &route
			}

		}

	}

	//Assume not public by default
	return nil
}

/**
Simple wrapping function that can be used else where.
*/
func SimpleLogger(inner http.Handler, name string) http.Handler {
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
