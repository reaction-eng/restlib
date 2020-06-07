// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/reaction-eng/restlib/users"
	"github.com/stretchr/testify/assert"
)

func TestHandler_handleOneTimePasswordGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                          string
		url                              string
		getUserUser                      func() users.User
		newUserUser                      func() users.User
		getUserByEmailCount              int
		getUserByEmailError              error
		issueOneTimePasswordRequestCount int
		issueOneTimePasswordRequestError error
		addUserCount                     int
		addUserError                     error
		addUserToOrgCount                int
		tokenGeneratorCount              int
		expectedStatus                   int
		expectedResponse                 string
	}{
		{
			comment: "working",
			url:     "http://localhost/users/onetimelogin?email=user@example.com&organizationId=43",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				return user
			},
			getUserByEmailCount:              1,
			getUserByEmailError:              nil,
			issueOneTimePasswordRequestCount: 1,
			issueOneTimePasswordRequestError: nil,
			newUserUser: func() users.User {
				return nil
			},
			addUserCount:        0,
			addUserError:        nil,
			tokenGeneratorCount: 1,
			addUserToOrgCount:   0,
			expectedStatus:      http.StatusOK,
			expectedResponse:    "{\"message\":\"onetimepassword_token_request_received\",\"status\":true}\n",
		},
		{
			comment: "user must specify email",
			url:     "http://localhost/users/onetimelogin?&organizationId=43",
			getUserUser: func() users.User {
				return nil
			},
			getUserByEmailCount:              0,
			getUserByEmailError:              nil,
			issueOneTimePasswordRequestCount: 0,
			issueOneTimePasswordRequestError: nil,
			newUserUser: func() users.User {
				return nil
			},
			addUserCount:        0,
			addUserError:        nil,
			tokenGeneratorCount: 0,
			expectedStatus:      http.StatusUnprocessableEntity,
			expectedResponse:    "{\"message\":\"onetimepassword_token_missing_email\",\"status\":false}\n",
		},
		{
			comment: "user must specify organizationId",
			url:     "http://localhost/users/onetimelogin?&email=user@example.com",
			getUserUser: func() users.User {
				return nil
			},
			getUserByEmailCount:              0,
			getUserByEmailError:              nil,
			issueOneTimePasswordRequestCount: 0,
			issueOneTimePasswordRequestError: nil,
			newUserUser: func() users.User {
				return nil
			},
			addUserCount:        0,
			addUserError:        nil,
			tokenGeneratorCount: 0,
			expectedStatus:      http.StatusUnprocessableEntity,
			expectedResponse:    "{\"message\":\"onetimepassword_token_missing_organizationId\",\"status\":false}\n",
		},
		{
			comment: "missing user should create a new user",
			url:     "http://localhost/users/onetimelogin?email=user@example.com&organizationId=43",
			getUserUser: func() users.User {
				return nil
			},
			getUserByEmailCount:              1,
			getUserByEmailError:              users.UserNotFound,
			issueOneTimePasswordRequestCount: 1,
			issueOneTimePasswordRequestError: nil,
			newUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetEmail("user@example.com").Times(1)
				user.EXPECT().SetPassword("").Times(1)
				user.EXPECT().SetOrganizations(43).Times(1)
				user.EXPECT().Organizations().Return([]int{43}).Times(1)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				return user
			},
			addUserCount:        1,
			addUserError:        nil,
			tokenGeneratorCount: 1,
			addUserToOrgCount:   1,
			expectedStatus:      http.StatusOK,
			expectedResponse:    "{\"message\":\"onetimepassword_token_request_received\",\"status\":true}\n",
		},
		{
			comment: "should error if add user db error",
			url:     "http://localhost/users/onetimelogin?email=user@example.com&organizationId=43",
			getUserUser: func() users.User {
				return nil
			},
			getUserByEmailCount: 1,
			getUserByEmailError: users.UserNotFound,
			newUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetEmail("user@example.com").Times(1)
				user.EXPECT().SetPassword("").Times(1)
				user.EXPECT().SetOrganizations(43).Times(1)
				user.EXPECT().Organizations().Return([]int{43}).Times(0)
				return user
			},
			addUserCount:                     1,
			addUserError:                     errors.New("db error"),
			issueOneTimePasswordRequestCount: 0,
			issueOneTimePasswordRequestError: nil,
			tokenGeneratorCount:              0,
			addUserToOrgCount:                0,
			expectedStatus:                   http.StatusForbidden,
			expectedResponse:                 "{\"message\":\"db error\",\"status\":false}\n",
		},
		{
			comment: "get user should error with other error",
			url:     "http://localhost/users/onetimelogin?email=user@example.com&organizationId=43",
			getUserUser: func() users.User {
				return nil
			},
			getUserByEmailCount:              1,
			getUserByEmailError:              errors.New("other error"),
			issueOneTimePasswordRequestCount: 0,
			issueOneTimePasswordRequestError: nil,
			newUserUser: func() users.User {
				return nil
			},
			addUserCount:        0,
			addUserError:        nil,
			tokenGeneratorCount: 0,
			expectedStatus:      http.StatusServiceUnavailable,
			expectedResponse:    "{\"message\":\"other error\",\"status\":false}\n",
		},
		{
			comment: "should return error is can't issue password ",
			url:     "http://localhost/users/onetimelogin?email=user@example.com&organizationId=43",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				return user
			},
			getUserByEmailCount: 1,
			getUserByEmailError: nil,
			newUserUser: func() users.User {
				return nil
			},
			addUserCount:                     0,
			addUserError:                     nil,
			issueOneTimePasswordRequestCount: 1,
			issueOneTimePasswordRequestError: errors.New("can't send email"),
			tokenGeneratorCount:              1,
			expectedStatus:                   http.StatusServiceUnavailable,
			expectedResponse:                 "{\"message\":\"can't send email\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockHelper.EXPECT().GetUserByEmail("user@example.com").Times(testCase.getUserByEmailCount).Return(testCase.getUserUser(), testCase.getUserByEmailError)
		token := "token123"

		if testCase.addUserCount > 0 {
			user := testCase.newUserUser()
			mockHelper.EXPECT().NewEmptyUser().Times(testCase.addUserCount).Return(user)
			mockHelper.EXPECT().AddUser(user).Times(testCase.addUserCount).Return(user, testCase.addUserError)
		}

		mockHelper.EXPECT().TokenGenerator().Times(testCase.tokenGeneratorCount).Return(token)
		mockHelper.EXPECT().IssueOneTimePasswordRequest(token, 34, "user@example.com").Times(testCase.issueOneTimePasswordRequestCount).Return(testCase.issueOneTimePasswordRequestError)
		mockHelper.EXPECT().AddUserToOrganization(gomock.Any(), 43).Times(testCase.addUserToOrgCount).Return(nil)

		handler := users.NewOneTimePasswordHandler(mockHelper)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("GET", testCase.url, nil)
		w := httptest.NewRecorder()

		// act
		route := router.Handle(w, req)

		// assert
		assert.NotNil(t, route)
		assert.Equal(t, true, route.Public)
		assert.Nil(t, route.ReqPermissions)
		assert.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, testCase.expectedResponse, w.Body.String())
	}
}

func TestHandler_handleOneTimePasswordLoginPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                           string
		body                              io.Reader
		getUserByEmailUser                func() users.User
		getUserByEmailCount               int
		getUserByEmailError               error
		checkForOneTimePasswordTokenCount int
		checkForOneTimePasswordTokenError error
		activateUserCount                 int
		activateUserError                 error
		getUserCount                      int
		getUserError                      error
		createJwtTokenCount               int
		useTokenCount                     int
		useTokenError                     error
		expectedStatus                    int
		expectedResponse                  string
	}{
		{
			comment: "working",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(2).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				user.EXPECT().Activated().Times(1).Return(true)
				user.EXPECT().Organizations().Times(1).Return([]int{45})
				user.EXPECT().SetToken("jwtToken").Times(1)

				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 1,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 0,
			activateUserError:                 nil,
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               1,
			useTokenCount:                     1,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusCreated,
			expectedResponse:                  "{}\n",
		},
		{
			comment: "mal formed json",
			body:    strings.NewReader(`{{}`),
			getUserByEmailUser: func() users.User {
				return nil
			},
			getUserByEmailCount:               0,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 0,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 0,
			activateUserError:                 nil,
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               0,
			useTokenCount:                     0,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusUnprocessableEntity,
			expectedResponse:                  "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
		{
			comment: "can't find user",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				return nil
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               users.UserNotFound,
			checkForOneTimePasswordTokenCount: 0,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 0,
			activateUserError:                 nil,
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               0,
			useTokenCount:                     0,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusForbidden,
			expectedResponse:                  "{\"message\":\"onetimepassword_forbidden\",\"status\":false}\n",
		},
		{
			comment: "user not in org",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Organizations().Times(1).Return([]int{1043})
				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 0,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 0,
			activateUserError:                 nil,
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               0,
			useTokenCount:                     0,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusForbidden,
			expectedResponse:                  "{\"message\":\"user_not_in_organization\",\"status\":false}\n",
		},
		{
			comment: "invalid token ",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Organizations().Times(1).Return([]int{45})

				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 1,
			checkForOneTimePasswordTokenError: errors.New("wrong token"),
			activateUserCount:                 0,
			activateUserError:                 nil,
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               0,
			useTokenCount:                     0,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusForbidden,
			expectedResponse:                  "{\"message\":\"onetimepassword_forbidden\",\"status\":false}\n",
		},
		{
			comment: "active users when need ",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(3).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				user.EXPECT().Activated().Times(1).Return(false)
				user.EXPECT().Organizations().Times(1).Return([]int{45})
				user.EXPECT().SetToken("jwtToken").Times(1)

				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 1,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 1,
			activateUserError:                 nil,
			getUserCount:                      1,
			getUserError:                      nil,
			createJwtTokenCount:               1,
			useTokenCount:                     1,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusCreated,
			expectedResponse:                  "{}\n",
		},
		{
			comment: "activate users when need returns if can't activate",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Activated().Times(1).Return(false)
				user.EXPECT().Organizations().Times(1).Return([]int{45})

				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 1,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 1,
			activateUserError:                 errors.New("can't activate"),
			getUserCount:                      0,
			getUserError:                      nil,
			createJwtTokenCount:               0,
			useTokenCount:                     1,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusServiceUnavailable,
			expectedResponse:                  "{\"message\":\"onetimepassword_forbidden\",\"status\":false}\n",
		},
		{
			comment: "activate users when need returns if can't get user",
			body:    strings.NewReader(`{"email":"user@example.com", "login_token":"login_token", "organizationId": 45 }`),
			getUserByEmailUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(2).Return(34)
				user.EXPECT().Activated().Times(1).Return(false)
				user.EXPECT().Organizations().Times(1).Return([]int{45})

				return user
			},
			getUserByEmailCount:               1,
			getUserByEmailError:               nil,
			checkForOneTimePasswordTokenCount: 1,
			checkForOneTimePasswordTokenError: nil,
			activateUserCount:                 1,
			activateUserError:                 nil,
			getUserCount:                      1,
			getUserError:                      errors.New("can't get user"),
			createJwtTokenCount:               0,
			useTokenCount:                     1,
			useTokenError:                     nil,
			expectedStatus:                    http.StatusServiceUnavailable,
			expectedResponse:                  "{\"message\":\"onetimepassword_forbidden\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		user := testCase.getUserByEmailUser()
		mockHelper.EXPECT().GetUserByEmail("user@example.com").Times(testCase.getUserByEmailCount).Return(user, testCase.getUserByEmailError)
		mockHelper.EXPECT().CheckForOneTimePasswordToken(34, "login_token").Times(testCase.checkForOneTimePasswordTokenCount).Return(103243, testCase.checkForOneTimePasswordTokenError)
		mockHelper.EXPECT().UseToken(103243).Times(testCase.useTokenCount).Return(testCase.useTokenError)

		mockHelper.EXPECT().ActivateUser(user).Times(testCase.activateUserCount).Return(testCase.activateUserError)
		mockHelper.EXPECT().GetUser(34).Times(testCase.getUserCount).Return(user, testCase.getUserError)

		mockHelper.EXPECT().CreateJWTToken(34, 45, "user@example.com").Times(testCase.createJwtTokenCount).Return("jwtToken")

		handler := users.NewOneTimePasswordHandler(mockHelper)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("POST", "http://localhost/users/onetimelogin", testCase.body)
		w := httptest.NewRecorder()

		// act
		route := router.Handle(w, req)

		// assert
		assert.NotNil(t, route)
		assert.Equal(t, true, route.Public)
		assert.Nil(t, route.ReqPermissions)
		assert.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, testCase.expectedResponse, w.Body.String())
	}
}
