package authentication

import (
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"bitbucket.org/reidev/restlib/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"strings"
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

			//check if request does not need authentication, serve the request if it doesn't need it
			if router.PublicRoute(requestPath) {
				//Just serve it
				next.ServeHTTP(w, r)
				return
			}

			//Get the header for auth
			tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

			//Token is missing, returns with error code 403 Unauthorized
			if tokenHeader == "" {
				//Return the error
				utils.ReturnJsonMessage(w, http.StatusForbidden, "Missing auth token")

				return
			}

			splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
			if len(splitted) != 2 {

				//Return the error
				utils.ReturnJsonMessage(w, http.StatusForbidden, "Invalid/Malformed auth token")

				return
			}

			//Grab the token part, what we are truly interested in
			tokenPart := splitted[1]

			//Get the token and take it back apart
			tk := &users.Token{}

			//Now parse the token
			token, err := jwt.ParseWithClaims(tokenPart, tk,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("token_password")), nil
				})

			//check for mailformed data
			if err != nil { //Malformed token, returns with http code 403 as usual
				//Return the error
				utils.ReturnJsonMessage(w, http.StatusForbidden, "Malformed authentication token")

				return
			}

			//Token is invalid, maybe not signed on this server
			if !token.Valid {
				//Return the error
				utils.ReturnJsonMessage(w, http.StatusForbidden, "Token is not valid.")

				return
			}

			//Now look up the user by id
			_, err = userRepo.GetUser(tk.UserId)

			//If there is an error return
			if err != nil {
				//Return the error
				utils.ReturnJsonError(w, http.StatusForbidden, err)

				return
			}
			//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
			//fmt.Sprintf("User %", tk.Username) //Useful for monitoring
			ctx := context.WithValue(r.Context(), "user", tk.UserId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r) //proceed in the middleware chain!
		})
	}
}
