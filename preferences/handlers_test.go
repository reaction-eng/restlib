package preferences_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		comment          string
		user             func() users.User
		userError        error
		preferences      *preferences.Preferences
		preferencesError error
		expectedStatus   int
		expectedResponse string
	}{
		{
			comment: "with user error",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			userError:        errors.New("user_db_error"),
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"user_db_error\",\"status\":false}\n",
		},
		{
			comment: "with preferences error",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusServiceUnavailable,
			expectedResponse: "{\"message\":\"preference_error\",\"status\":false}\n",
		},
		{
			comment: "with user and preferences errors",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			userError:        errors.New("user_db_error"),
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"user_db_error\",\"status\":false}\n",
		},
		{
			comment: "all valid",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			preferences:      &preferences.Preferences{Settings: &preferences.SettingGroup{Settings: map[string]string{"info": "123"}}},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"settings\":{\"settings\":{\"info\":\"123\"},\"subgroup\":null},\"options\":null}\n",
		},
		{
			comment:          "no user logged in",
			user:             func() users.User { return nil },
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().GetUser(101).MaxTimes(1).Return(testCase.user(), testCase.userError)

		mockPrefRepo := mocks.NewMockPreferencesRepo(mockCtrl)
		mockPrefRepo.EXPECT().GetPreferences(testCase.user()).MaxTimes(1).Return(testCase.preferences, testCase.preferencesError)

		handler := preferences.NewHandler(mockUserRepo, mockPrefRepo)

		router := mocks.NewTestRouterWithUser(handler, testCase.user(), 0)

		req := httptest.NewRequest("GET", "http://localhost/users/preferences", nil)
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

func TestHandler_handleUserPreferencesSet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment          string
		user             func() users.User
		userError        error
		body             io.Reader
		preferences      *preferences.Preferences
		preferencesError error
		expectedStatus   int
		expectedResponse string
	}{
		{
			comment:          "no user logged in",
			user:             func() users.User { return nil },
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"no_user_logged_in\",\"status\":false}\n",
		},

		{
			comment: "with user error",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			userError:        errors.New("user_db_error"),
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"user_db_error\",\"status\":false}\n",
		},
		{
			comment: "with user and preference error",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			userError:        errors.New("user_db_error"),
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"user_db_error\",\"status\":false}\n",
		},
		{
			comment: "with preference error and body",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			body:             strings.NewReader("{\"info\":\"123\"}"),
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusServiceUnavailable,
			expectedResponse: "{\"message\":\"preference_error\",\"status\":false}\n",
		},
		{
			comment: "without body",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"unexpected end of JSON input\",\"status\":false}\n",
		},
		{
			comment: "with bad body",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			body:             strings.NewReader("{\"settings\"null}"),
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"invalid character 'n' after object key\",\"status\":false}\n",
		},
		{
			comment: "with user and preference error and body",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			body:             strings.NewReader("{\"info\":\"123\"}"),
			userError:        errors.New("user_db_error"),
			preferencesError: errors.New("preference_error"),
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "{\"message\":\"user_db_error\",\"status\":false}\n",
		},
		{
			comment: "all good",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(101).MaxTimes(1)
				return user
			},
			body:             strings.NewReader("{\"info\":\"123\"}"),
			preferences:      &preferences.Preferences{Settings: &preferences.SettingGroup{Settings: map[string]string{"info": "123"}}},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"settings\":{\"settings\":{\"info\":\"123\"},\"subgroup\":null},\"options\":null}\n",
		},
	}

	// arrange
	for _, testCase := range testCases {
		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().GetUser(101).MaxTimes(1).Return(testCase.user(), testCase.userError)

		mockPrefRepo := mocks.NewMockPreferencesRepo(mockCtrl)
		mockPrefRepo.EXPECT().SetPreferences(testCase.user(), gomock.Any()).MaxTimes(1).Return(testCase.preferences, testCase.preferencesError)

		handler := preferences.NewHandler(mockUserRepo, mockPrefRepo)

		router := mocks.NewTestRouterWithUser(handler, testCase.user(), 0)

		req := httptest.NewRequest("POST", "http://localhost/users/preferences", testCase.body)
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
