// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reaction-eng/restlib/roles"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/reaction-eng/restlib/users"
)

func TestHandler_handleUserPermissionsGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                string
		getUserUser            func() users.User
		getUserCount           int
		getUserError           error
		getPermissionsCount    int
		getPermissionsResponse []string
		getPermissionsError    error
		expectedStatus         int
		expectedResponse       string
	}{
		{
			comment: "working",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(1023)
				return user
			},
			getUserCount:           1,
			getUserError:           nil,
			getPermissionsCount:    1,
			getPermissionsResponse: []string{"perm1", "perm2"},
			getPermissionsError:    nil,
			expectedStatus:         http.StatusOK,
			expectedResponse:       "{\"permissions\":[\"perm1\",\"perm2\"]}\n",
		},
		{
			comment: "not logged in",
			getUserUser: func() users.User {
				return nil
			},
			getUserCount:           0,
			getUserError:           nil,
			getPermissionsCount:    0,
			getPermissionsResponse: []string{},
			getPermissionsError:    nil,
			expectedStatus:         http.StatusForbidden,
			expectedResponse:       "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
		{
			comment: "can't find user with db error",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(1023)
				return user
			},
			getUserCount:           1,
			getUserError:           errors.New("can't find user"),
			getPermissionsCount:    0,
			getPermissionsResponse: []string{"perm1", "perm2"},
			getPermissionsError:    nil,
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"can't find user\",\"status\":false}\n",
		},
		{
			comment: "can't find user",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(1023)
				return user
			},
			getUserCount:           1,
			getUserError:           users.UserNotFound,
			getPermissionsCount:    0,
			getPermissionsResponse: []string{"perm1", "perm2"},
			getPermissionsError:    nil,
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"" + users.UserNotFound.Error() + "\",\"status\":false}\n",
		},
		{
			comment: "can't find permissions",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(1023)
				return user
			},
			getUserCount:           1,
			getUserError:           nil,
			getPermissionsCount:    1,
			getPermissionsResponse: []string{"perm1", "perm2"},
			getPermissionsError:    errors.New("can't get info"),
			expectedStatus:         http.StatusServiceUnavailable,
			expectedResponse:       "{\"message\":\"can't get info\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		user := testCase.getUserUser()
		mockUserRepo.EXPECT().GetUser(1023).Times(testCase.getUserCount).Return(user, testCase.getUserError)

		permissions := roles.Permissions{
			Permissions: testCase.getPermissionsResponse,
		}

		mockRoleRepo := mocks.NewMockRolesRepo(mockCtrl)
		mockRoleRepo.EXPECT().GetPermissions(user, 1000).Times(testCase.getPermissionsCount).Return(&permissions, testCase.getPermissionsError)

		handler := roles.NewHandler(mockUserRepo, mockRoleRepo)
		router := mocks.NewTestRouterWithUser(handler, user, 1000)

		req := httptest.NewRequest("GET", "http://localhost/users/permissions", nil)
		w := httptest.NewRecorder()

		// act
		route := router.Handle(w, req)

		// assert
		assert.NotNil(t, route)
		assert.Equal(t, false, route.Public)
		assert.Nil(t, route.ReqPermissions)
		assert.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, testCase.expectedResponse, w.Body.String())
	}
}
