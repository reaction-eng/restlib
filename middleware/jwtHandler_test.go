// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reaction-eng/restlib/roles"
	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"

	"github.com/golang/mock/gomock"

	"github.com/reaction-eng/restlib/mocks"

	"github.com/reaction-eng/restlib/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMakeJwtMiddlewareFunc(t *testing.T) {
	testCases := []struct {
		description string
		method      string
		header      map[string]string
		setup       func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper)
		status      int
		response    string
		next        bool
		context     map[string]interface{}
	}{
		{
			"options",
			"OPTIONS",
			map[string]string{},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {

			},
			http.StatusOK,
			"",
			true,
			nil,
		},
		{
			"get without known route",
			"GET",
			map[string]string{},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(nil)

			},
			http.StatusForbidden,
			"{\"message\":\"\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get without token",
			"GET",
			map[string]string{},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)
			},
			http.StatusForbidden,
			"{\"message\":\"auth_missing_token\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with public route",
			"GET",
			map[string]string{},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: true}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

			},
			http.StatusOK,
			"",
			true,
			nil,
		},
		{
			"get with public route with token",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: true}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

			},
			http.StatusOK,
			"",
			true,
			nil,
		},
		{
			"get with private route with bad token",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(0, 0, "", errors.New("bad token"))

			},
			http.StatusForbidden,
			"{\"message\":\"bad token\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with bad websocket token",
			"GET",
			map[string]string{"Sec-Websocket-Protocol": "Example_Space_one two three,"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example one two three").Return(0, 0, "", errors.New("bad token"))

			},
			http.StatusForbidden,
			"{\"message\":\"bad token\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but no user",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "", nil)

				mockUsers.EXPECT().GetUser(100).Times(1).Return(nil, nil)

			},
			http.StatusForbidden,
			"{\"message\":\"auth_malformed_token\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but user error",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "", nil)

				mockUser := mocks.NewMockUser(ctrl)

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, errors.New("user error"))

			},
			http.StatusForbidden,
			"{\"message\":\"user error\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but valid user with wrong email",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("")

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

			},
			http.StatusForbidden,
			"{\"message\":\"auth_malformed_token\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but non activated user",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{Public: false}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(false)

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

			},
			http.StatusForbidden,
			"{\"message\":\"user_not_activated\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but without permission",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{"req_permission"},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{1000})

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(1).Return(&roles.Permissions{
					Permissions: []string{},
				}, nil)
			},
			http.StatusForbidden,
			"{\"message\":\"insufficient_access\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but without permission",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{"req_permission"},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{1000})

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(1).Return(&roles.Permissions{
					Permissions: []string{"req_permissio"},
				}, nil)
			},
			http.StatusForbidden,
			"{\"message\":\"insufficient_access\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token but with permission error",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{1000})

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(1).Return(&roles.Permissions{
					Permissions: []string{},
				}, errors.New("permission error"))
			},
			http.StatusForbidden,
			"{\"message\":\"insufficient_access\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get with private route with good token wrong org",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{132})

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(0).Return(&roles.Permissions{
					Permissions: []string{},
				}, errors.New("user_not_in_organization"))
			},
			http.StatusForbidden,
			"{\"message\":\"user_not_in_organization\",\"status\":false}\n",
			false,
			nil,
		},
		{
			"get user context with private route without needed permissions",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{1000})

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(1).Return(&roles.Permissions{
					Permissions: []string{},
				}, nil)
			},
			http.StatusOK,
			"",
			true,
			map[string]interface{}{utils.UserKey: 100, utils.OrganizationKey: 1000},
		},
		{
			"get user context with private route with needed permissions",
			"GET",
			map[string]string{"Authorization": "Example Token"},
			func(ctrl *gomock.Controller, mockRouter *mocks.MockRouter, mockUsers *mocks.MockUserRepo, mockRoles *mocks.MockRolesRepo, mockHelper *mocks.MockHelper) {
				route := &routing.Route{
					Public:         false,
					ReqPermissions: []string{"req_perm"},
				}

				mockRouter.EXPECT().GetRoute(gomock.Any()).Times(1).Return(route)

				mockHelper.EXPECT().ValidateToken("Example Token").Return(100, 1000, "example@example.com", nil)

				mockUser := mocks.NewMockUser(ctrl)
				mockUser.EXPECT().Email().Times(1).Return("example@example.com")
				mockUser.EXPECT().Activated().Times(1).Return(true)
				mockUser.EXPECT().Organizations().Times(1).Return([]int{1000})

				mockUsers.EXPECT().GetUser(100).Times(1).Return(mockUser, nil)

				mockRoles.EXPECT().GetPermissions(mockUser, 1000).Times(1).Return(&roles.Permissions{
					Permissions: []string{"perm 1", "perm 2", "req_perm", "perm 4"},
				}, nil)
			},
			http.StatusOK,
			"",
			true,
			map[string]interface{}{utils.UserKey: 100, utils.OrganizationKey: 1000},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest(testCase.method, "http://localhost/example", nil)

		var wResponse http.ResponseWriter
		var rResponse *http.Request
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wResponse = w
			rResponse = r
		})

		for k, v := range testCase.header {
			r.Header.Set(k, v)
		}

		mockRouter := mocks.NewMockRouter(mockCtrl)
		mockUsers := mocks.NewMockUserRepo(mockCtrl)
		mockRoles := mocks.NewMockRolesRepo(mockCtrl)
		mockHelper := mocks.NewMockHelper(mockCtrl)

		testCase.setup(mockCtrl, mockRouter, mockUsers, mockRoles, mockHelper)

		middleware := middleware.MakeJwtMiddlewareFunc(mockRouter, mockUsers, mockRoles, mockHelper)
		handler := middleware.Middleware(mockHandler)

		// act
		handler.ServeHTTP(w, r)

		// assert
		if testCase.next {
			for k, v := range testCase.context {
				value := rResponse.Context().Value(k)
				assert.Equal(t, v, value, testCase.description)
			}

			assert.Equal(t, w, wResponse, testCase.description)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode, testCase.description)
			assert.Empty(t, w.Body.String(), testCase.description)
		} else {
			assert.Equal(t, testCase.status, w.Result().StatusCode, testCase.description)
			assert.Equal(t, testCase.response, w.Body.String(), testCase.description)
			assert.Nil(t, wResponse, "should not proceed to the next handler", testCase.description)
			assert.Nil(t, rResponse, "should not proceed to the next handler", testCase.description)
		}
	}
}
