// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

/**
Define a function to allow CORS
Allow CORS here By * or specific origin
*/
func MakeCORSMiddlewareFunc() mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin,authorization,content-type, x-ijt, X-Auth-Token, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")

			//Just serve it
			next.ServeHTTP(w, r)
		})
	}
}
