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

func TestNewHandler(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	mockUserHelper := mocks.NewMockUserHelper(mockCtrl)

	// act
	handler := users.NewHandler(mockUserHelper, true)

	// assert
	assert.NotNil(t, handler)
}

func TestHandler_handleUserCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                string
		body                   io.Reader
		expectedUser           func() users.User
		expectDecodeSuccessful bool
		createUserError        error
		expectedStatus         int
		expectedResponse       string
	}{
		{
			comment: "working",
			body:    strings.NewReader(`{"email":"user@example.info", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetEmail("user@example.info").Times(1)
				user.EXPECT().SetPassword("new password").Times(1)
				user.EXPECT().SetOrganizations(34).Times(1)
				return user
			},
			expectDecodeSuccessful: true,
			createUserError:        nil,
			expectedStatus:         http.StatusCreated,
			expectedResponse:       "{\"message\":\"create_user_added\",\"status\":true}\n",
		},
		{
			comment: "bad json",
			body:    strings.NewReader(`{{"email":"user@example.info", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			expectDecodeSuccessful: false,
			createUserError:        nil,
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
		{
			comment: "can't add user",
			body:    strings.NewReader(`{"email":"user@example.info", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetEmail("user@example.info").Times(1)
				user.EXPECT().SetPassword("new password").Times(1)
				user.EXPECT().SetOrganizations(34).Times(1)
				return user
			},
			expectDecodeSuccessful: true,
			createUserError:        errors.New("can't add user"),
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"can't add user\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockUser := testCase.expectedUser()
		mockHelper.EXPECT().NewEmptyUser().Times(1).Return(mockUser)

		if testCase.expectDecodeSuccessful {
			mockHelper.EXPECT().CreateUser(mockUser).Times(1).Return(testCase.createUserError)
		}

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("POST", "http://localhost/users/new", testCase.body)
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

func TestHandler_handleUserLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                     string
		body                        io.Reader
		expectedUser                func() users.User
		expectedPassword            string
		expectedOrgId               int
		expectedGetUserByEmailCount int
		getUserByEmailError         error
		loginCount                  int
		loginError                  error
		expectedStatus              int
		expectedResponse            string
	}{
		{
			comment: "working",
			body:    strings.NewReader(`{"email":"user@example.info", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			getUserByEmailError:         nil,
			expectedGetUserByEmailCount: 1,
			loginError:                  nil,
			loginCount:                  1,
			expectedStatus:              http.StatusCreated,
			expectedResponse:            "{}\n", // empty user return
		},
		{
			comment: "should clean up the email name",
			body:    strings.NewReader(`{"email":" uSeR@example.info ", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			getUserByEmailError:         nil,
			expectedGetUserByEmailCount: 1,
			loginError:                  nil,
			loginCount:                  1,
			expectedStatus:              http.StatusCreated,
			expectedResponse:            "{}\n", // empty user return
		},
		{
			comment: "can't find user error",
			body:    strings.NewReader(`{"email":" uSeR@example.info ", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			getUserByEmailError:         errors.New("can't find user"),
			expectedGetUserByEmailCount: 1,
			loginError:                  nil,
			loginCount:                  0,
			expectedStatus:              http.StatusForbidden,
			expectedResponse:            "{\"message\":\"can't find user\",\"status\":false}\n",
		},
		{
			comment: "can't login error",
			body:    strings.NewReader(`{"email":" uSeR@example.info ", "password":"new password", "organizationId": 34 }`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			getUserByEmailError:         nil,
			expectedGetUserByEmailCount: 1,
			loginError:                  errors.New("wrong password"),
			loginCount:                  1,
			expectedStatus:              http.StatusForbidden,
			expectedResponse:            "{\"message\":\"wrong password\",\"status\":false}\n",
		},
		{
			comment: "can't decode error",
			body:    strings.NewReader(`{{}`),
			expectedUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				return user
			},
			getUserByEmailError:         nil,
			expectedGetUserByEmailCount: 0,
			loginError:                  nil,
			loginCount:                  0,
			expectedStatus:              http.StatusUnprocessableEntity,
			expectedResponse:            "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockUser := testCase.expectedUser()

		mockHelper.EXPECT().GetUserByEmail("user@example.info").Times(testCase.expectedGetUserByEmailCount).Return(mockUser, testCase.getUserByEmailError)
		mockHelper.EXPECT().Login("new password", 34, mockUser).Times(testCase.loginCount).Return(mockUser, testCase.loginError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("POST", "http://localhost/users/login", testCase.body)
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

func TestHandler_handleUserUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type testUser struct {
		users.BasicUser
		OtherInfo string `json:"otherInfo"`
	}

	loggedInUser := mocks.NewMockUser(mockCtrl)
	loggedInUser.EXPECT().Id().MinTimes(0).Return(34)

	testCases := []struct {
		comment                 string
		body                    io.Reader
		loggedInUser            users.User
		expectedUserAfterDecode func() users.User
		getUserError            error
		getUserCount            int
		updateError             error
		updateCount             int
		expectedStatus          int
		expectedResponse        string
	}{
		{
			comment:      "nobody logged in",
			body:         strings.NewReader(`{"email":"user@example.info", "otherInfo":"new other info" }`),
			loggedInUser: nil,
			expectedUserAfterDecode: func() users.User {
				user := &testUser{}
				return user
			},
			getUserError:     nil,
			getUserCount:     0,
			updateError:      nil,
			updateCount:      0,
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
		{
			comment:      "working",
			body:         strings.NewReader(`{"email":"user@example.info", "otherInfo":"new other info" }`),
			loggedInUser: loggedInUser,
			expectedUserAfterDecode: func() users.User {
				user := &testUser{
					OtherInfo: "new other info",
					BasicUser: users.BasicUser{
						Id_:    34,
						Email_: "user@example.info",
					},
				}
				return user
			},
			getUserError:     nil,
			getUserCount:     1,
			updateError:      nil,
			updateCount:      1,
			expectedStatus:   http.StatusAccepted,
			expectedResponse: "{\"id\":34,\"organizations\":null,\"email\":\"user@example.info\",\"token\":\"\",\"otherInfo\":\"new other info\"}\n", // empty user return
		},
		{
			comment:      "can't find user",
			body:         strings.NewReader(`{"email":"user@example.info", "otherInfo":"new other info" }`),
			loggedInUser: loggedInUser,
			expectedUserAfterDecode: func() users.User {
				user := &testUser{}
				return user
			},
			getUserError:     errors.New("can't get user"),
			getUserCount:     1,
			updateError:      nil,
			updateCount:      0,
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"can't get user\",\"status\":false}\n",
		},
		{
			comment:      "can't decode",
			body:         strings.NewReader(`{{}`),
			loggedInUser: loggedInUser,
			expectedUserAfterDecode: func() users.User {
				user := &testUser{}
				return user
			},
			getUserError:     nil,
			getUserCount:     1,
			updateError:      nil,
			updateCount:      0,
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
		{
			comment:      "can't update user user",
			body:         strings.NewReader(`{"email":"user@example.info", "otherInfo":"new other info" }`),
			loggedInUser: loggedInUser,
			expectedUserAfterDecode: func() users.User {
				user := &testUser{
					OtherInfo: "new other info",
					BasicUser: users.BasicUser{
						Id_:    34,
						Email_: "user@example.info",
					},
				}
				return user
			},
			getUserError:     nil,
			getUserCount:     1,
			updateError:      errors.New("can't update user"),
			updateCount:      1,
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: "{\"message\":\"can't update user\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		user := &testUser{
			BasicUser: users.BasicUser{
				Id_: 34,
			},
		}

		mockHelper.EXPECT().GetUser(34).Times(testCase.getUserCount).Return(user, testCase.getUserError)
		mockHelper.EXPECT().Update(34, testCase.expectedUserAfterDecode()).Times(testCase.updateCount).Return(testCase.expectedUserAfterDecode(), testCase.updateError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouterWithUser(handler, testCase.loggedInUser, 0) // not logged in

		req := httptest.NewRequest("PUT", "http://localhost/users/", testCase.body)
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

func TestHandler_handleUserGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	loggedInUser := mocks.NewMockUser(mockCtrl)
	loggedInUser.EXPECT().Id().MinTimes(0).Return(34)

	testCases := []struct {
		comment          string
		loggedInUser     users.User
		getUserUser      func() users.User
		getUserError     error
		getUserCount     int
		expectedStatus   int
		expectedResponse string
	}{
		{
			comment:      "nobody logged in",
			loggedInUser: nil,
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetPassword("").Times(0)
				return user
			},
			getUserError:     nil,
			getUserCount:     0,
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
		{
			comment:      "everything working",
			loggedInUser: loggedInUser,
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().SetPassword("").Times(1)
				return user
			},
			getUserError:     nil,
			getUserCount:     1,
			expectedStatus:   http.StatusOK,
			expectedResponse: "{}\n",
		},
		{
			comment:      "everything working",
			loggedInUser: loggedInUser,
			getUserUser: func() users.User {
				return nil
			},
			getUserError:     errors.New("can't get user"),
			getUserCount:     1,
			expectedStatus:   http.StatusUnsupportedMediaType,
			expectedResponse: "{\"message\":\"can't get user\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockHelper.EXPECT().GetUser(34).Times(testCase.getUserCount).Return(testCase.getUserUser(), testCase.getUserError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouterWithUser(handler, testCase.loggedInUser, 0)

		req := httptest.NewRequest("GET", "http://localhost/users/", nil)
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

func TestHandler_handlePasswordUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	loggedInUser := mocks.NewMockUser(mockCtrl)
	loggedInUser.EXPECT().Id().MinTimes(0).Return(34)

	testCases := []struct {
		comment             string
		body                io.Reader
		loggedInUser        users.User
		passwordChangeCount int
		passwordChangeError error
		expectedStatus      int
		expectedResponse    string
	}{
		{
			comment:             "working",
			body:                strings.NewReader(`{"email":"user@example.info", "password":"new password", "passwordold": "old password" }`),
			loggedInUser:        loggedInUser,
			passwordChangeCount: 1,
			passwordChangeError: nil,
			expectedStatus:      http.StatusAccepted,
			expectedResponse:    "{\"message\":\"password_change_success\",\"status\":true}\n",
		},
		{
			comment:             "nobody logged in",
			body:                strings.NewReader(`{"email":"user@example.info", "password":"new password", "passwordold": "old password" }`),
			loggedInUser:        nil,
			passwordChangeCount: 0,
			passwordChangeError: nil,
			expectedStatus:      http.StatusForbidden,
			expectedResponse:    "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
		{
			comment:             "PasswordChangeError",
			body:                strings.NewReader(`{"email":"user@example.info", "password":"new password", "passwordold": "old password" }`),
			loggedInUser:        loggedInUser,
			passwordChangeCount: 1,
			passwordChangeError: errors.New("password change error"),
			expectedStatus:      http.StatusForbidden,
			expectedResponse:    "{\"message\":\"password change error\",\"status\":false}\n",
		},
		{
			comment:             "Decode Error",
			body:                strings.NewReader(`{{}`),
			loggedInUser:        loggedInUser,
			passwordChangeCount: 0,
			passwordChangeError: nil,
			expectedStatus:      http.StatusUnprocessableEntity,
			expectedResponse:    "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockHelper.EXPECT().PasswordChange(34, "user@example.info", "new password", "old password").Times(testCase.passwordChangeCount).Return(testCase.passwordChangeError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouterWithUser(handler, testCase.loggedInUser, 0)

		req := httptest.NewRequest("POST", "http://localhost/users/password/change", testCase.body)
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

func TestHandler_handlePasswordResetGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                string
		url                    string
		getUserUser            func() users.User
		getUserByEmailCount    int
		getUserByEmailError    error
		issueResetRequestCount int
		issueResetRequestError error
		tokenGeneratorCount    int
		expectedStatus         int
		expectedResponse       string
	}{
		{
			comment: "working",
			url:     "http://localhost/users/password/reset?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				return user
			},
			getUserByEmailCount:    1,
			getUserByEmailError:    nil,
			issueResetRequestCount: 1,
			issueResetRequestError: nil,
			tokenGeneratorCount:    1,
			expectedStatus:         http.StatusOK,
			expectedResponse:       "{\"message\":\"password_change_request_received\",\"status\":true}\n",
		},
		{
			comment: "no email specified",
			url:     "http://localhost/users/password/reset",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				return user
			},
			getUserByEmailCount:    0,
			getUserByEmailError:    nil,
			issueResetRequestCount: 0,
			issueResetRequestError: nil,
			tokenGeneratorCount:    0,
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"password_change_missing_email\",\"status\":false}\n",
		},
		{
			comment: "empty email",
			url:     "http://localhost/users/password/reset?email=",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				return user
			},
			getUserByEmailCount:    0,
			getUserByEmailError:    nil,
			issueResetRequestCount: 0,
			issueResetRequestError: nil,
			tokenGeneratorCount:    0,
			expectedStatus:         http.StatusUnprocessableEntity,
			expectedResponse:       "{\"message\":\"password_change_missing_email\",\"status\":false}\n",
		},
		{
			comment: "can't find user",
			url:     "http://localhost/users/password/reset?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				return user
			},
			getUserByEmailCount:    1,
			getUserByEmailError:    errors.New("can't find user"),
			issueResetRequestCount: 0,
			issueResetRequestError: nil,
			tokenGeneratorCount:    0,
			expectedStatus:         http.StatusOK,
			expectedResponse:       "{\"message\":\"password_change_request_received\",\"status\":true}\n",
		},
		{
			comment: "can't issue request",
			url:     "http://localhost/users/password/reset?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				return user
			},
			getUserByEmailCount:    1,
			getUserByEmailError:    nil,
			issueResetRequestCount: 1,
			issueResetRequestError: errors.New("can't issue request"),
			tokenGeneratorCount:    1,
			expectedStatus:         http.StatusServiceUnavailable,
			expectedResponse:       "{\"message\":\"can't issue request\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockHelper.EXPECT().GetUserByEmail("user@example.com").Times(testCase.getUserByEmailCount).Return(testCase.getUserUser(), testCase.getUserByEmailError)
		token := "token123"
		mockHelper.EXPECT().TokenGenerator().Times(testCase.tokenGeneratorCount).Return(token)
		mockHelper.EXPECT().IssueResetRequest(token, 34, "user@example.com").Times(testCase.issueResetRequestCount).Return(testCase.issueResetRequestError)

		handler := users.NewHandler(mockHelper, true)
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

func TestHandler_handlePasswordResetPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                   string
		body                      io.Reader
		getUserUser               func() users.User
		getUserByEmailCount       int
		getUserByEmailError       error
		checkForResetTokenCount   int
		checkForResetTokenError   error
		passwordChangeForcedCount int
		passwordChangeForcedError error
		useTokenCount             int
		userTokenError            error
		expectedStatus            int
		expectedResponse          string
	}{
		{
			comment: "working",
			body:    strings.NewReader(`{"email":"user@example.info", "reset_token":"reset 123", "password": "new password" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(2).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       1,
			getUserByEmailError:       nil,
			checkForResetTokenCount:   1,
			checkForResetTokenError:   nil,
			passwordChangeForcedCount: 1,
			passwordChangeForcedError: nil,
			useTokenCount:             1,
			userTokenError:            nil,
			expectedStatus:            http.StatusAccepted,
			expectedResponse:          "{\"message\":\"password_change_success\",\"status\":true}\n",
		},
		{
			comment: "bad decode",
			body:    strings.NewReader(`{{}`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       0,
			getUserByEmailError:       nil,
			checkForResetTokenCount:   0,
			checkForResetTokenError:   nil,
			passwordChangeForcedCount: 0,
			passwordChangeForcedError: nil,
			useTokenCount:             0,
			userTokenError:            nil,
			expectedStatus:            http.StatusUnprocessableEntity,
			expectedResponse:          "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
		{
			comment: "no user",
			body:    strings.NewReader(`{"email":"user@example.info", "reset_token":"reset 123", "password": "new password" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       1,
			getUserByEmailError:       errors.New("no user error"),
			checkForResetTokenCount:   0,
			checkForResetTokenError:   nil,
			passwordChangeForcedCount: 0,
			passwordChangeForcedError: nil,
			useTokenCount:             0,
			userTokenError:            nil,
			expectedStatus:            http.StatusForbidden,
			expectedResponse:          "{\"message\":\"password_change_forbidden\",\"status\":false}\n",
		},
		{
			comment: "no token error",
			body:    strings.NewReader(`{"email":"user@example.info", "reset_token":"reset 123", "password": "new password" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       1,
			getUserByEmailError:       nil,
			checkForResetTokenCount:   1,
			checkForResetTokenError:   errors.New("no token error"),
			passwordChangeForcedCount: 0,
			passwordChangeForcedError: nil,
			useTokenCount:             0,
			userTokenError:            nil,
			expectedStatus:            http.StatusForbidden,
			expectedResponse:          "{\"message\":\"password_change_forbidden\",\"status\":false}\n",
		},
		{
			comment: "can't change password",
			body:    strings.NewReader(`{"email":"user@example.info", "reset_token":"reset 123", "password": "new password" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(2).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       1,
			getUserByEmailError:       nil,
			checkForResetTokenCount:   1,
			checkForResetTokenError:   nil,
			passwordChangeForcedCount: 1,
			passwordChangeForcedError: errors.New("can't change password"),
			useTokenCount:             0,
			userTokenError:            nil,
			expectedStatus:            http.StatusForbidden,
			expectedResponse:          "{\"message\":\"can't change password\",\"status\":false}\n",
		},
		{
			comment: "user token error",
			body:    strings.NewReader(`{"email":"user@example.info", "reset_token":"reset 123", "password": "new password" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(2).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				return user
			},
			getUserByEmailCount:       1,
			getUserByEmailError:       nil,
			checkForResetTokenCount:   1,
			checkForResetTokenError:   nil,
			passwordChangeForcedCount: 1,
			passwordChangeForcedError: nil,
			useTokenCount:             1,
			userTokenError:            errors.New("user token error"),
			expectedStatus:            http.StatusForbidden,
			expectedResponse:          "{\"message\":\"user token error\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		reqId := 454

		mockHelper.EXPECT().GetUserByEmail("user@example.info").Times(testCase.getUserByEmailCount).Return(testCase.getUserUser(), testCase.getUserByEmailError)
		mockHelper.EXPECT().CheckForResetToken(34, "reset 123").Times(testCase.checkForResetTokenCount).Return(reqId, testCase.checkForResetTokenError)
		mockHelper.EXPECT().PasswordChangeForced(34, "user@example.info", "new password").Times(testCase.passwordChangeForcedCount).Return(testCase.passwordChangeForcedError)
		mockHelper.EXPECT().UseToken(reqId).Times(testCase.useTokenCount).Return(testCase.userTokenError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("POST", "http://localhost/users/password/reset", testCase.body)
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

func TestHandler_handleUserActivationPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                      string
		body                         io.Reader
		getUserUser                  func() users.User
		getUserByEmailCount          int
		getUserByEmailError          error
		checkForActivationTokenCount int
		checkForActivationError      error
		activateUserCount            int
		activateUserError            error
		useTokenCount                int
		userTokenError               error
		expectedStatus               int
		expectedResponse             string
	}{
		{
			comment: "working",
			body:    strings.NewReader(`{"email":"user@example.info", "activation_token":"token 123" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				return user
			},
			getUserByEmailCount:          1,
			getUserByEmailError:          nil,
			checkForActivationTokenCount: 1,
			checkForActivationError:      nil,
			activateUserCount:            1,
			activateUserError:            nil,
			useTokenCount:                1,
			userTokenError:               nil,
			expectedStatus:               http.StatusAccepted,
			expectedResponse:             "{\"message\":\"user_activated\",\"status\":true}\n",
		},
		{
			comment: "decode error",
			body:    strings.NewReader(`{{}}`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				return user
			},
			getUserByEmailCount:          0,
			getUserByEmailError:          nil,
			checkForActivationTokenCount: 0,
			checkForActivationError:      nil,
			activateUserCount:            0,
			activateUserError:            nil,
			useTokenCount:                0,
			userTokenError:               nil,
			expectedStatus:               http.StatusUnprocessableEntity,
			expectedResponse:             "{\"message\":\"invalid character '{' looking for beginning of object key string\",\"status\":false}\n",
		},
		{
			comment: "can't find error",
			body:    strings.NewReader(`{"email":"user@example.info", "activation_token":"token 123" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				return user
			},
			getUserByEmailCount:          1,
			getUserByEmailError:          errors.New("can't find user"),
			checkForActivationTokenCount: 0,
			checkForActivationError:      nil,
			activateUserCount:            0,
			activateUserError:            nil,
			useTokenCount:                0,
			userTokenError:               nil,
			expectedStatus:               http.StatusForbidden,
			expectedResponse:             "{\"message\":\"activation_forbidden\",\"status\":false}\n",
		},
		{
			comment: "bad token",
			body:    strings.NewReader(`{"email":"user@example.info", "activation_token":"token 123" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				return user
			},
			getUserByEmailCount:          1,
			getUserByEmailError:          nil,
			checkForActivationTokenCount: 1,
			checkForActivationError:      errors.New("bad token"),
			activateUserCount:            0,
			activateUserError:            nil,
			useTokenCount:                0,
			userTokenError:               nil,
			expectedStatus:               http.StatusForbidden,
			expectedResponse:             "{\"message\":\"activation_forbidden\",\"status\":false}\n",
		},
		{
			comment: "can't activate",
			body:    strings.NewReader(`{"email":"user@example.info", "activation_token":"token 123" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				return user
			},
			getUserByEmailCount:          1,
			getUserByEmailError:          nil,
			checkForActivationTokenCount: 1,
			checkForActivationError:      nil,
			activateUserCount:            1,
			activateUserError:            errors.New("can't activate user"),
			useTokenCount:                0,
			userTokenError:               nil,
			expectedStatus:               http.StatusForbidden,
			expectedResponse:             "{\"message\":\"can't activate user\",\"status\":false}\n",
		},
		{
			comment: "can't use token",
			body:    strings.NewReader(`{"email":"user@example.info", "activation_token":"token 123" }`),
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				return user
			},
			getUserByEmailCount:          1,
			getUserByEmailError:          nil,
			checkForActivationTokenCount: 1,
			checkForActivationError:      nil,
			activateUserCount:            1,
			activateUserError:            nil,
			useTokenCount:                1,
			userTokenError:               errors.New("can't use token"),
			expectedStatus:               http.StatusForbidden,
			expectedResponse:             "{\"message\":\"can't use token\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		reqId := 454

		user := testCase.getUserUser()

		mockHelper.EXPECT().GetUserByEmail("user@example.info").Times(testCase.getUserByEmailCount).Return(user, testCase.getUserByEmailError)
		mockHelper.EXPECT().CheckForActivationToken(34, "token 123").Times(testCase.checkForActivationTokenCount).Return(reqId, testCase.checkForActivationError)
		mockHelper.EXPECT().ActivateUser(user).Times(testCase.activateUserCount).Return(testCase.activateUserError)
		mockHelper.EXPECT().UseToken(reqId).Times(testCase.useTokenCount).Return(testCase.userTokenError)

		handler := users.NewHandler(mockHelper, true)
		router := mocks.NewTestRouter(handler) // not logged in

		req := httptest.NewRequest("POST", "http://localhost/users/activate", testCase.body)
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

func TestHandler_handleUserActivationGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                     string
		url                         string
		getUserUser                 func() users.User
		getUserByEmailCount         int
		getUserByEmailError         error
		issueActivationRequestCount int
		issueActivationRequestError error
		tokenGeneratorCount         int
		expectedStatus              int
		expectedResponse            string
	}{
		{
			comment: "working",
			url:     "http://localhost/users/activate?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				user.EXPECT().Activated().Times(1).Return(false)
				return user
			},
			getUserByEmailCount:         1,
			getUserByEmailError:         nil,
			issueActivationRequestCount: 1,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         1,
			expectedStatus:              http.StatusOK,
			expectedResponse:            "{\"message\":\"activation_token_request_received\",\"status\":true}\n",
		},
		{
			comment: "no email specified",
			url:     "http://localhost/users/activate",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				user.EXPECT().Activated().Times(0).Return(false)
				return user
			},
			getUserByEmailCount:         0,
			getUserByEmailError:         nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			expectedStatus:              http.StatusUnprocessableEntity,
			expectedResponse:            "{\"message\":\"activation_token_missing_email\",\"status\":false}\n",
		},
		{
			comment: "empty email",
			url:     "http://localhost/users/activate?email=",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				user.EXPECT().Activated().Times(0).Return(false)
				return user
			},
			getUserByEmailCount:         0,
			getUserByEmailError:         nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			expectedStatus:              http.StatusUnprocessableEntity,
			expectedResponse:            "{\"message\":\"activation_token_missing_email\",\"status\":false}\n",
		},
		{
			comment: "can't find user",
			url:     "http://localhost/users/activate?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				user.EXPECT().Activated().Times(0).Return(false)
				return user
			},
			getUserByEmailCount:         1,
			getUserByEmailError:         errors.New("can't find user"),
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			expectedStatus:              http.StatusOK,
			expectedResponse:            "{\"message\":\"activation_token_request_received\",\"status\":true}\n",
		},
		{
			comment: "can't issue request",
			url:     "http://localhost/users/activate?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.com")
				user.EXPECT().Activated().Times(1).Return(false)
				return user
			},
			getUserByEmailCount:         1,
			getUserByEmailError:         nil,
			issueActivationRequestCount: 1,
			issueActivationRequestError: errors.New("can't issue request"),
			tokenGeneratorCount:         1,
			expectedStatus:              http.StatusServiceUnavailable,
			expectedResponse:            "{\"message\":\"can't issue request\",\"status\":false}\n",
		},
		{
			comment: "don't double issue request",
			url:     "http://localhost/users/activate?email=user@example.com",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.com")
				user.EXPECT().Activated().Times(1).Return(true)
				return user
			},
			getUserByEmailCount:         1,
			getUserByEmailError:         nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			expectedStatus:              http.StatusOK,
			expectedResponse:            "{\"message\":\"activation_token_request_received\",\"status\":true}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockHelper := mocks.NewMockUserHelper(mockCtrl)

		mockHelper.EXPECT().GetUserByEmail("user@example.com").Times(testCase.getUserByEmailCount).Return(testCase.getUserUser(), testCase.getUserByEmailError)
		token := "token123"
		mockHelper.EXPECT().TokenGenerator().Times(testCase.tokenGeneratorCount).Return(token)
		mockHelper.EXPECT().IssueActivationRequest(token, 34, "user@example.com").Times(testCase.issueActivationRequestCount).Return(testCase.issueActivationRequestError)

		handler := users.NewHandler(mockHelper, true)
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
