package users_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/reaction-eng/restlib/users"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHelper(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
	mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)
	mockPasswordHelper := mocks.NewMockHelper(mockCtrl)

	// act
	basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

	// assert
	assert.NotNil(t, basicHelper)
	assert.Equal(t, mockUserRepo, basicHelper.Repo)
	assert.Equal(t, mockPasswordResetRepo, basicHelper.ResetRepo)
	assert.Equal(t, mockPasswordHelper, basicHelper.Helper)
}

func TestBasicHelper_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                     string
		user                        func() users.User
		hashPasswordCount           int
		addUserCount                int
		addUserError                error
		validatePasswordCount       int
		validatePasswordError       error
		issueActivationRequestCount int
		issueActivationRequestError error
		tokenGeneratorCount         int
		addUserToOrgCount           int
		expectedError               error
	}{
		{
			comment: "working",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(2).Return("user@example.info")
				user.EXPECT().Password().Times(2).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(1)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(1)
				return user
			},
			validatePasswordCount:       1,
			validatePasswordError:       nil,
			hashPasswordCount:           1,
			addUserCount:                1,
			addUserError:                nil,
			issueActivationRequestCount: 1,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         1,
			addUserToOrgCount:           2,
			expectedError:               nil,
		},
		{
			comment: "invalid password",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(0)
				return user
			},
			validatePasswordCount:       1,
			validatePasswordError:       errors.New("invalid password"),
			hashPasswordCount:           0,
			addUserCount:                0,
			addUserError:                nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			addUserToOrgCount:           0,
			expectedError:               errors.New("invalid password"),
		},
		{
			comment: "not valid email",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(1).Return("userexample.info")
				user.EXPECT().Password().Times(0).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(0)
				return user
			},
			validatePasswordCount:       0,
			validatePasswordError:       nil,
			hashPasswordCount:           0,
			addUserCount:                0,
			addUserError:                nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			addUserToOrgCount:           0,
			expectedError:               errors.New("validate_missing_email"),
		},
		{
			comment: "empty email",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(1).Return("")
				user.EXPECT().Password().Times(0).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(0)
				return user
			},
			validatePasswordCount:       0,
			validatePasswordError:       nil,
			hashPasswordCount:           0,
			addUserCount:                0,
			addUserError:                nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			addUserToOrgCount:           0,
			expectedError:               errors.New("validate_missing_email"),
		},
		{
			comment: "all white space",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(1).Return("  			 ")
				user.EXPECT().Password().Times(0).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(0)
				return user
			},
			validatePasswordCount:       0,
			validatePasswordError:       nil,
			hashPasswordCount:           0,
			addUserCount:                0,
			addUserError:                nil,
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			addUserToOrgCount:           0,
			expectedError:               errors.New("validate_missing_email"),
		},
		{
			comment: "add user error, user already there",
			user: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(2).Return("password 123")
				user.EXPECT().SetPassword("hashed password").Times(1)
				user.EXPECT().Organizations().Return([]int{1000, 1002}).Times(0)
				return user
			},
			validatePasswordCount:       1,
			validatePasswordError:       nil,
			hashPasswordCount:           1,
			addUserCount:                1,
			addUserError:                errors.New("user already here"),
			issueActivationRequestCount: 0,
			issueActivationRequestError: nil,
			tokenGeneratorCount:         0,
			addUserToOrgCount:           0,
			expectedError:               errors.New("user already here"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		user := testCase.user()

		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().AddUser(user).Times(testCase.addUserCount).Return(user, testCase.addUserError)
		mockUserRepo.EXPECT().AddUserToOrganization(user, gomock.Any()).Times(testCase.addUserToOrgCount).Return(nil)

		mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)
		token := "token 123"
		mockPasswordResetRepo.EXPECT().IssueActivationRequest(token, 34, "user@example.info").Times(testCase.issueActivationRequestCount).Return(testCase.issueActivationRequestError)

		mockPasswordHelper := mocks.NewMockHelper(mockCtrl)
		mockPasswordHelper.EXPECT().TokenGenerator().Times(testCase.tokenGeneratorCount).Return(token)
		mockPasswordHelper.EXPECT().HashPassword("password 123").Times(testCase.hashPasswordCount).Return("hashed password")
		mockPasswordHelper.EXPECT().ValidatePassword("password 123").Times(testCase.validatePasswordCount).Return(testCase.validatePasswordError)

		basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

		// act
		err := basicHelper.CreateUser(user)

		// assert
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestBasicHelper_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type userStruct struct {
		id       int
		password string
		email    string
		orgs     []int
	}

	testCases := []struct {
		comment         string
		newUser         userStruct
		oldUser         userStruct
		getUserError    error
		getUserCount    int
		updateUserError error
		updateUserCount int
		expectedError   error
		expectedUser    bool
	}{
		{
			comment: "working",
			newUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 1,
			expectedError:   nil,
			expectedUser:    true,
		},
		{
			comment: "emails don't match",
			newUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.comm",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 0,
			expectedError:   errors.New("update_forbidden"),
			expectedUser:    false,
		},
		{
			comment: "ids don't match",
			newUser: userStruct{
				id:       33,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 0,
			expectedError:   errors.New("update_forbidden"),
			expectedUser:    false,
		},
		{
			comment: "passwords don't match",
			newUser: userStruct{
				id:       34,
				password: "hashed password new",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 0,
			expectedError:   errors.New("update_forbidden"),
			expectedUser:    false,
		},
		{
			comment: "orgs don't match",
			newUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{43, 23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 0,
			expectedError:   errors.New("update_forbidden"),
			expectedUser:    false,
		},
		{
			comment: "get user error",
			newUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    errors.New("get user error"),
			getUserCount:    1,
			updateUserError: nil,
			updateUserCount: 0,
			expectedError:   errors.New("get user error"),
			expectedUser:    false,
		},
		{
			comment: "update user error",
			newUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			oldUser: userStruct{
				id:       34,
				password: "hashed password",
				email:    "user@example.com",
				orgs:     []int{23, 543},
			},
			getUserError:    nil,
			getUserCount:    1,
			updateUserError: errors.New("update user error"),
			updateUserCount: 1,
			expectedError:   errors.New("update user error"),
			expectedUser:    true,
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockNewUser := mocks.NewMockUser(mockCtrl)
		mockNewUser.EXPECT().Id().MaxTimes(1).Return(testCase.newUser.id)
		mockNewUser.EXPECT().Password().MaxTimes(1).Return(testCase.newUser.password)
		mockNewUser.EXPECT().Email().MaxTimes(1).Return(testCase.newUser.email)
		mockNewUser.EXPECT().Organizations().MaxTimes(1).Return(testCase.newUser.orgs)

		mockOldUser := mocks.NewMockUser(mockCtrl)
		mockOldUser.EXPECT().Id().MaxTimes(1).Return(testCase.oldUser.id)
		mockOldUser.EXPECT().Password().MaxTimes(1).Return(testCase.oldUser.password)
		mockOldUser.EXPECT().Email().MaxTimes(1).Return(testCase.oldUser.email)
		mockOldUser.EXPECT().Organizations().MaxTimes(1).Return(testCase.oldUser.orgs)

		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().GetUser(34).Times(testCase.getUserCount).Return(mockOldUser, testCase.getUserError)
		mockUserRepo.EXPECT().UpdateUser(mockNewUser).Times(testCase.updateUserCount).Return(mockNewUser, testCase.updateUserError)

		mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)
		mockPasswordHelper := mocks.NewMockHelper(mockCtrl)

		basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

		// act
		user, err := basicHelper.Update(34, mockNewUser)

		// assert
		if testCase.expectedUser {
			assert.Equal(t, mockNewUser, user)
		} else {
			assert.Nil(t, user)
		}
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestBasicHelper_PasswordChange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                  string
		userId                   int
		email                    string
		newPassword              string
		oldPassword              string
		getUserUser              func() users.User
		getUserError             error
		comparePasswordsCount    int
		comparePasswordsResponse bool
		validatePasswordCount    int
		validatePasswordError    error
		hashPasswordCount        int
		hashPasswordResponse     string
		updateUserCount          int
		updateUserError          error
		expectedError            error
	}{
		{
			comment:     "working",
			userId:      34,
			email:       "user@example.info",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(1)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: true,
			validatePasswordCount:    1,
			validatePasswordError:    nil,
			hashPasswordCount:        1,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          1,
			updateUserError:          nil,
			expectedError:            nil,
		},
		{
			comment:     "working with different case email",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(1)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: true,
			validatePasswordCount:    1,
			validatePasswordError:    nil,
			hashPasswordCount:        1,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          1,
			updateUserError:          nil,
			expectedError:            nil,
		},
		{
			comment:     "no user",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().PasswordLogin().Times(0).Return(true)
				return user
			},
			getUserError:             errors.New("missing user"),
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			hashPasswordCount:        0,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          0,
			updateUserError:          nil,
			expectedError:            errors.New("missing user"),
		},
		{
			comment:     "no password access",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().PasswordLogin().Times(1).Return(false)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			hashPasswordCount:        0,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          0,
			updateUserError:          nil,
			expectedError:            errors.New("user_password_login_forbidden"),
		}, {
			comment:     "emails don't match",
			userId:      34,
			email:       " uSeR@example.inFoo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			hashPasswordCount:        0,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          0,
			updateUserError:          nil,
			expectedError:            errors.New("password_change_forbidden"),
		},
		{
			comment:     "wrong password",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: false,
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			hashPasswordCount:        0,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          0,
			updateUserError:          nil,
			expectedError:            errors.New("password_change_forbidden"),
		},
		{
			comment:     "invalid password",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(0)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: true,
			validatePasswordCount:    1,
			validatePasswordError:    errors.New("bad password"),
			hashPasswordCount:        0,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          0,
			updateUserError:          nil,
			expectedError:            errors.New("bad password"),
		},
		{
			comment:     "can't update user",
			userId:      34,
			email:       " uSeR@example.inFo ",
			newPassword: "new password",
			oldPassword: "old password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("hashed password").Times(1)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				return user
			},
			getUserError:             nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: true,
			validatePasswordCount:    1,
			validatePasswordError:    nil,
			hashPasswordCount:        1,
			hashPasswordResponse:     "hashed password",
			updateUserCount:          1,
			updateUserError:          errors.New("can't update user"),
			expectedError:            errors.New("can't update user"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		user := testCase.getUserUser()

		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().GetUser(testCase.userId).Times(1).Return(user, testCase.getUserError)
		mockUserRepo.EXPECT().UpdateUser(user).Times(testCase.updateUserCount).Return(nil, testCase.updateUserError)

		mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)

		mockPasswordHelper := mocks.NewMockHelper(mockCtrl)
		mockPasswordHelper.EXPECT().ComparePasswords("hashed password", testCase.oldPassword).Times(testCase.comparePasswordsCount).Return(testCase.comparePasswordsResponse)
		mockPasswordHelper.EXPECT().ValidatePassword("new password").Times(testCase.validatePasswordCount).Return(testCase.validatePasswordError)
		mockPasswordHelper.EXPECT().HashPassword("new password").Times(testCase.hashPasswordCount).Return("hashed password")

		basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

		// act
		err := basicHelper.PasswordChange(testCase.userId, testCase.email, testCase.newPassword, testCase.oldPassword)

		// assert
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestBasicHelper_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                  string
		getUserUser              func() users.User
		validatePasswordCount    int
		validatePasswordError    error
		comparePasswordsCount    int
		comparePasswordsResponse bool
		createJWTTokenCount      int
		expectedError            error
	}{
		{
			comment: "working",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				user.EXPECT().Activated().Times(1).Return(true)
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("").Times(1)
				user.EXPECT().Id().Times(1).Return(34)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(1)
				user.EXPECT().Organizations().Times(1).Return([]int{3, 454, 534})
				return user
			},
			validatePasswordCount:    1,
			validatePasswordError:    nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: true,
			createJWTTokenCount:      1,
			expectedError:            nil,
		},
		{
			comment: "wrong org",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				user.EXPECT().Activated().Times(1).Return(true)
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("").Times(0)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(0)
				user.EXPECT().Organizations().Times(1).Return([]int{3, 2, 534})
				return user
			},
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			createJWTTokenCount:      0,
			expectedError:            errors.New("user_not_in_organization"),
		},
		{
			comment: "can't use a password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(false)
				user.EXPECT().Activated().Times(0).Return(true)
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("").Times(0)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(0)
				user.EXPECT().Organizations().Times(0).Return([]int{3, 454, 534})

				return user
			},
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			createJWTTokenCount:      0,
			expectedError:            errors.New("user_password_login_forbidden"),
		},
		{
			comment: "not activated",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				user.EXPECT().Activated().Times(1).Return(false)
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("").Times(0)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(0)
				user.EXPECT().Organizations().Times(0).Return([]int{3, 454, 534})

				return user
			},
			validatePasswordCount:    0,
			validatePasswordError:    nil,
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			createJWTTokenCount:      0,
			expectedError:            errors.New("user_not_activated"),
		},
		{
			comment: "non valid password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				user.EXPECT().Activated().Times(1).Return(true)
				user.EXPECT().Password().Times(0).Return("hashed password")
				user.EXPECT().SetPassword("").Times(0)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(0)
				user.EXPECT().Organizations().Times(1).Return([]int{3, 454, 534})
				return user
			},
			validatePasswordCount:    1,
			validatePasswordError:    errors.New("bad password"),
			comparePasswordsCount:    0,
			comparePasswordsResponse: true,
			createJWTTokenCount:      0,
			expectedError:            errors.New("login_invalid_password"),
		},
		{
			comment: "pass words don't match",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().PasswordLogin().Times(1).Return(true)
				user.EXPECT().Activated().Times(1).Return(true)
				user.EXPECT().Password().Times(1).Return("hashed password")
				user.EXPECT().SetPassword("").Times(1)
				user.EXPECT().Id().Times(0).Return(34)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetToken("token").Times(0)
				user.EXPECT().Organizations().Times(1).Return([]int{3, 454, 534})
				return user
			},
			validatePasswordCount:    1,
			validatePasswordError:    nil,
			comparePasswordsCount:    1,
			comparePasswordsResponse: false,
			createJWTTokenCount:      0,
			expectedError:            errors.New("login_invalid_password"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		user := testCase.getUserUser()

		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)

		mockPasswordHelper := mocks.NewMockHelper(mockCtrl)
		mockPasswordHelper.EXPECT().ValidatePassword("new password").Times(testCase.validatePasswordCount).Return(testCase.validatePasswordError)
		mockPasswordHelper.EXPECT().ComparePasswords("hashed password", "new password").Times(testCase.comparePasswordsCount).Return(testCase.comparePasswordsResponse)
		mockPasswordHelper.EXPECT().CreateJWTToken(34, 454, "user@example.info").Times(testCase.createJWTTokenCount).Return("token")

		basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

		// act
		returnUser, err := basicHelper.Login("new password", 454, user)

		// assert
		if testCase.expectedError == nil {
			assert.Equal(t, user, returnUser)
		} else {
			assert.Nil(t, returnUser)
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestBasicHelper_PasswordChangeForced(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment               string
		userId                int
		email                 string
		newPassword           string
		getUserUser           func() users.User
		getUserError          error
		validatePasswordCount int
		validatePasswordError error
		hashPasswordCount     int
		hashPasswordResponse  string
		updateUserCount       int
		updateUserError       error
		expectedError         error
	}{
		{
			comment:     "working",
			userId:      34,
			email:       "user@example.info",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(1)
				return user
			},
			getUserError:          nil,
			validatePasswordCount: 1,
			validatePasswordError: nil,
			hashPasswordCount:     1,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       1,
			updateUserError:       nil,
			expectedError:         nil,
		},
		{
			comment: "weird email",
			userId:  34,
			email: " 	uSeR@example.INfo		",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(1)
				return user
			},
			getUserError:          nil,
			validatePasswordCount: 1,
			validatePasswordError: nil,
			hashPasswordCount:     1,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       1,
			updateUserError:       nil,
			expectedError:         nil,
		},
		{
			comment: "can't find user",
			userId:  34,
			email: " 	uSeR@example.INfo		",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(0).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(0)
				return user
			},
			getUserError:          errors.New("can't find user"),
			validatePasswordCount: 0,
			validatePasswordError: nil,
			hashPasswordCount:     0,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       0,
			updateUserError:       nil,
			expectedError:         errors.New("can't find user"),
		},
		{
			comment:     "emails don't match",
			userId:      34,
			email:       "user@example.infoo",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(0)
				return user
			},
			getUserError:          nil,
			validatePasswordCount: 0,
			validatePasswordError: nil,
			hashPasswordCount:     0,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       0,
			updateUserError:       nil,
			expectedError:         errors.New("password_change_forbidden"),
		},
		{
			comment:     "non valid passsword",
			userId:      34,
			email:       "user@example.info",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(0)
				return user
			},
			getUserError:          nil,
			validatePasswordCount: 1,
			validatePasswordError: errors.New("non valid password"),
			hashPasswordCount:     0,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       0,
			updateUserError:       nil,
			expectedError:         errors.New("non valid password"),
		},
		{
			comment:     "can't update user",
			userId:      34,
			email:       "user@example.info",
			newPassword: "new password",
			getUserUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Times(1).Return("user@example.info")
				user.EXPECT().SetPassword("hashed password").Times(1)
				return user
			},
			getUserError:          nil,
			validatePasswordCount: 1,
			validatePasswordError: nil,
			hashPasswordCount:     1,
			hashPasswordResponse:  "hashed password",
			updateUserCount:       1,
			updateUserError:       errors.New("can't update user"),
			expectedError:         errors.New("can't update user"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		user := testCase.getUserUser()

		mockUserRepo := mocks.NewMockUserRepo(mockCtrl)
		mockUserRepo.EXPECT().GetUser(testCase.userId).Times(1).Return(user, testCase.getUserError)
		mockUserRepo.EXPECT().UpdateUser(user).Times(testCase.updateUserCount).Return(nil, testCase.updateUserError)

		mockPasswordResetRepo := mocks.NewMockResetRepo(mockCtrl)

		mockPasswordHelper := mocks.NewMockHelper(mockCtrl)
		mockPasswordHelper.EXPECT().ValidatePassword("new password").Times(testCase.validatePasswordCount).Return(testCase.validatePasswordError)
		mockPasswordHelper.EXPECT().HashPassword("new password").Times(testCase.hashPasswordCount).Return("hashed password")

		basicHelper := users.NewUserHelper(mockUserRepo, mockPasswordResetRepo, mockPasswordHelper)

		// act
		err := basicHelper.PasswordChangeForced(testCase.userId, testCase.email, testCase.newPassword)

		// assert
		assert.Equal(t, testCase.expectedError, err)
	}
}
