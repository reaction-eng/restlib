package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reaction-eng/restlib/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHerokuHttpsOnlyMiddlewareFunc(t *testing.T) {
	testCases := []struct {
		headerKey   string
		headerValue string
		redirect    bool
	}{
		{"x-forwarded-proto", "https", false},
		{"x-forwarded-proto", "http", true},
		{"x-forwarded-proto", "something", true},
		{"other header", "https", true},
		{"*", "*", true},
	}

	for _, testCase := range testCases {
		// arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/example", nil)

		var wResponse http.ResponseWriter
		var rResponse *http.Request
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wResponse = w
			rResponse = r
		})

		r.Header.Set(testCase.headerKey, testCase.headerValue)

		middleware := middleware.HerokuHttpsOnlyMiddlewareFunc()
		handler := middleware.Middleware(mockHandler)

		// act
		handler.ServeHTTP(w, r)

		// assert
		if testCase.redirect {
			assert.Equal(t, http.StatusTemporaryRedirect, w.Result().StatusCode)
			assert.Equal(t, "<a href=\"https://localhosthttp://localhost/example\">Temporary Redirect</a>.\n\n", w.Body.String())
			assert.Nil(t, wResponse)
			assert.Nil(t, rResponse)
		} else {
			assert.Equal(t, w, wResponse)
			assert.Equal(t, r, rResponse)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Empty(t, w.Body.String())
		}
	}
}
