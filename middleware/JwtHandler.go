// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package middleware

import (
	"github.com/reaction-eng/restlib/passwords"
	"github.com/reaction-eng/restlib/roles"
	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/users"
	"github.com/reaction-eng/restlib/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

/**
Define a function to handle checking for auth
*/
func MakeJwtMiddlewareFunc(router *routing.Router, userRepo users.Repo, permRepo roles.Repo, passHelper passwords.Helper) mux.MiddlewareFunc {

	//Return an instance
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//If this is options just bypass
			if r.Method == "OPTIONS" {
				//Just continue to serve as usual
				next.ServeHTTP(w, r)
				return
			}

			//current the current route
			route := router.GetRoute(r)

			//If the route was not found return
			if route == nil {
				//Return the error
				utils.ReturnJsonStatus(w, http.StatusForbidden, false, "")
				return
			}

			//tokenHeader will get set here if we have a websocket. If this isn't a websocket, tokenHeader will be ""
			//var tokenHeader string
			tokenHeader := r.Header.Get("Sec-Websocket-Protocol")

			//if true, it's not a websocket. If it has something, we are dealing with a websocket.
			if tokenHeader == "" {

				//check if request does not need middleware, serve the request if it doesn't need it
				if route.Public {
					//Just serve it
					next.ServeHTTP(w, r)
					return
				}
				//Get the header for auth
				tokenHeader = r.Header.Get("Authorization") //Grab the token from the header
			} else {
				tokenHeader = strings.Replace(tokenHeader, "_Space_", " ", -1)
				locOfComma := strings.Index(tokenHeader, ",")
				tokenHeader = tokenHeader[0:locOfComma]
			}

			//Validate and get the user id
			userId, tokenEmail, err := passHelper.ValidateToken(tokenHeader)

			//If there is an error return
			if err != nil {
				//Return the error
				utils.ReturnJsonError(w, http.StatusForbidden, err)

				return
			}

			//Now look up the user by id
			loggedInUser, err := userRepo.GetUser(userId)

			//If there is an error return
			if err != nil {
				//Return the error
				utils.ReturnJsonError(w, http.StatusForbidden, err)

				return
			}
			//Make sure the emails match in the token and logged in user
			if loggedInUser.Email() != tokenEmail {
				//Return the error
				utils.ReturnJsonStatus(w, http.StatusForbidden, false, "auth_malformed_token")

				return
			}

			//Make sure that the person is activated
			if !loggedInUser.Activated() {
				//There prob is not a user to return
				utils.ReturnJsonStatus(w, http.StatusForbidden, false, "user_not_activated")
			}

			//Make sure that the user has permission
			if permRepo != nil {
				//See if we are allowed
				userPerm, err := permRepo.GetPermissions(loggedInUser)

				//See if we are allowed to
				if err != nil || !userPerm.AllowedTo(route.ReqPermissions...) {
					//Return the error
					utils.ReturnJsonStatus(w, http.StatusForbidden, false, "insufficient_access")
					return
				}

			}

			//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
			//fmt.Sprintf("User %", tk.Username) //Useful for monitoring
			ctx := context.WithValue(r.Context(), "user", userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r) //proceed in the middleware chain!
		})
	}
}
