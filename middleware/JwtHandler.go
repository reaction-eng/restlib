package middleware

import (
	"bitbucket.org/reidev/restlib/passwords"
	"bitbucket.org/reidev/restlib/roles"
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"bitbucket.org/reidev/restlib/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
)

/**
Define a function to handle checking for auth
*/
func MakeJwtMiddlewareFunc(router *routing.Router, userRepo users.Repo, permRepo roles.Repo) mux.MiddlewareFunc {

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

			//check if request does not need middleware, serve the request if it doesn't need it
			if route.Public {
				//Just serve it
				next.ServeHTTP(w, r)
				return
			}

			//Get the header for auth
			tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

			//Validate and get the user id
			userId, tokenEmail, err := passwords.ValidateToken(tokenHeader)

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
