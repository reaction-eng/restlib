// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reaction-eng/restlib/users"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/reaction-eng/restlib/preferences"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
	mockPrefRepo := mocks.NewMockPreferencesRepo(mockCtrl)

	// act
	handler := preferences.NewHandler(mockUserRepo, mockPrefRepo)

	// assert
	assert.NotNil(t, handler)
}

func TestHandler_handleUserPreferencesGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		user     func() users.User
		status   int
		response string
	}{
		/*{
			func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).Times(1)

				return user
			},
		},*/
		{
			func() users.User { return nil },
			http.StatusForbidden,
			"{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockPrefRepo := mocks.NewMockPreferencesRepo(mockCtrl)

		handler := preferences.NewHandler(mockUserRepo, mockPrefRepo)

		router := mocks.NewTestRouter(handler, testCase.user())

		req := httptest.NewRequest("GET", "http://localhost/users/preferences", nil)
		w := httptest.NewRecorder()

		// act
		route := router.Handle(w, req)

		// assert
		assert.NotNil(t, route)
		assert.Equal(t, false, route.Public)
		assert.Nil(t, route.ReqPermissions)
		assert.Equal(t, testCase.status, w.Result().StatusCode)
		assert.Equal(t, testCase.response, w.Body.String())
	}

}
