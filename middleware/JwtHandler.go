package middleware

import (
	"bitbucket.org/reidev/restlib/authentication"
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
func MakeJwtMiddlewareFunc(router *routing.Router, userRepo users.Repo) mux.MiddlewareFunc {

	//Return an instance
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//If this is options just bypass
			if r.Method == "OPTIONS" {
				//Just continue to serve as usual
				next.ServeHTTP(w, r)
				return
			}

			//current request path
			requestPath := r.URL.Path

			//check if request does not need middleware, serve the request if it doesn't need it
			if router.PublicRoute(requestPath) {
				//Just serve it
				next.ServeHTTP(w, r)
				return
			}

			//Get the header for auth
			tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

			//Validate and get the user id
			userId, err := authentication.ValidateToken(tokenHeader)

			//If there is an error return
			if err != nil {
				//Return the error
				utils.ReturnJsonError(w, http.StatusForbidden, err)

				return
			}

			//Now look up the user by id
			_, err = userRepo.GetUser(userId)

			//If there is an error return
			if err != nil {
				//Return the error
				utils.ReturnJsonError(w, http.StatusForbidden, err)

				return
			}
			//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
			//fmt.Sprintf("User %", tk.Username) //Useful for monitoring
			ctx := context.WithValue(r.Context(), "user", userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r) //proceed in the middleware chain!
		})
	}
}
