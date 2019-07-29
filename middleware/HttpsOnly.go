package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

/**
Define a function to allow CORS
Allow CORS here By * or specific origin
*/
func HerokuHttpsOnlyMiddlewareFunc() mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
