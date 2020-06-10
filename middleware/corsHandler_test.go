// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reaction-eng/restlib/middleware"

	"github.com/stretchr/testify/assert"
)

func TestMakeCORSMiddlewareFunc(t *testing.T) {
	// arrange
	r := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()

	var wResponse http.ResponseWriter
	var rResponse *http.Request
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wResponse = w
		rResponse = r
	})

	middleware := middleware.MakeCORSMiddlewareFunc()
	handler := middleware.Middleware(mockHandler)

	// act
	handler.ServeHTTP(w, r)

	// assert
	assert.Equal(t, w, wResponse)
	assert.Equal(t, r, rResponse)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin,authorization,content-type, x-ijt, X-Auth-Token, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "GET, POST, PATCH, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
}
