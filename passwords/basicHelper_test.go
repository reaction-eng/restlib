// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/reaction-eng/restlib/passwords"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicHelper(t *testing.T) {
	testCases := []struct {
		configuration map[string]string
		expectedError error
		expectResult  bool
	}{
		{
			configuration: map[string]string{"token_password": "UOGWZSRAODMTCMYZOUFXUOGWZSRAODMTCMYZOUFXUOGWZSRAODMTCMYZOUFX"},
			expectedError: nil,
			expectResult:  true,
		},
		{
			configuration: map[string]string{"token_password": "UOGWZSRAODMTCZSRAODMTCMYZOUFX"},
			expectedError: errors.New("the jwt token is not specified or not long enough"),
			expectResult:  false,
		},
		{
			configuration: map[string]string{"token_password": ""},
			expectedError: errors.New("the jwt token is not specified or not long enough"),
			expectResult:  false,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
		for k, v := range testCase.configuration {
			mockConfiguration.EXPECT().GetString(k).Return(v).Times(1)
		}

		// act
		helper, err := passwords.NewBasicHelper(mockConfiguration)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if testCase.expectResult {
			assert.NotNil(t, helper)
		} else {
			assert.Nil(t, helper)
		}
	}
}

func setupBasicHelper(t *testing.T, mockCtrl *gomock.Controller) *passwords.BasicHelper {
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetString("token_password").Return("UOGWZSRAODMTCMYZOUFXUOGWZSRAODMTCMYZOUFXUOGWZSRAODMTCMYZOUFX").Times(1)

	helper, err := passwords.NewBasicHelper(mockConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	return helper
}

func TestBasicHelper_HashPassword(t *testing.T) {
	testCases := []struct {
		input          string
		expectedLength int
	}{
		{
			"123456",
			60,
		},
		{
			"alphaBetaGama",
			60,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		helper := setupBasicHelper(t, mockCtrl)

		// act
		output := helper.HashPassword(testCase.input)

		// assert
		assert.Equal(t, testCase.expectedLength, len(output))
	}
}

func TestBasicHelper_CreateJWTToken(t *testing.T) {
	testCases := []struct {
		userId         int
		orgId          int
		email          string
		expectedResult string
	}{
		{
			42,
			65,
			"example@example.com",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQyLCJPcmdhbml6YXRpb25JZCI6NjUsIkVtYWlsIjoiZXhhbXBsZUBleGFtcGxlLmNvbSJ9.8xvP_tiVqrowq85_t4eoJqh0PXJGqOY1mG6ixKcFpqw",
		},
		{
			102,
			23,
			"example2@example.com",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEwMiwiT3JnYW5pemF0aW9uSWQiOjIzLCJFbWFpbCI6ImV4YW1wbGUyQGV4YW1wbGUuY29tIn0.s0IMktQm5I3DJ9Bixdyd42-q3dEbj6xAFz_v7AOaLRY",
		},
		{
			98,
			346,
			"matt@example.com",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjk4LCJPcmdhbml6YXRpb25JZCI6MzQ2LCJFbWFpbCI6Im1hdHRAZXhhbXBsZS5jb20ifQ.pzf0sT2F77gH-1Ghqw_bwWrxSDbpc85FD1nXxBByHQc",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		helper := setupBasicHelper(t, mockCtrl)

		// act
		output := helper.CreateJWTToken(testCase.userId, testCase.orgId, testCase.email)

		// assert
		assert.Equal(t, testCase.expectedResult, output)
	}
}

func TestBasicHelper_ComparePasswords(t *testing.T) {
	testCases := []struct {
		currentPwHash string
		testPassword  string
		equals        bool
	}{
		{
			"$2y$10$Dh6eCPAO43w2Ta72Mqu./.YDjMjXFyNd3Jl4AAyqukC/iDblorUNu",
			"123456",
			true,
		},
		{
			"$2y$10$Dh6eCPAO43w2Ta72Mqu./.YDjMjXFyNd3Jl4AAyqukC/iDblorUNu",
			"12345",
			false,
		},
		{
			"$2y$10$mgBz16gBtT1UodLYcapX7ORIZR7Il3bwulDqz5fGQHqCIh2WQ5No2",
			"alphaBetaGama",
			true,
		},
		{
			"$2y$10$mgBz16gBtT1UodLYcapX7ORIZR7Il3bwulDqz5fGQHqCIh2WQ5No2",
			"alphafBetaGama",
			false,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		helper := setupBasicHelper(t, mockCtrl)

		// act
		equals := helper.ComparePasswords(testCase.currentPwHash, testCase.testPassword)

		// assert
		assert.Equal(t, testCase.equals, equals)
	}
}

func TestBasicHelper_TokenGenerator(t *testing.T) {
	// arrange
	testRounds := 1000

	previous := make([]string, 0)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	helper := setupBasicHelper(t, mockCtrl)

	// act
	// assert
	for i := 0; i < testRounds; i++ {
		newToken := helper.TokenGenerator()

		for _, value := range previous {
			assert.NotEqual(t, newToken, value)
		}
		previous = append(previous, newToken)
	}
}

func TestBasicHelper_ValidateToken(t *testing.T) {
	testCases := []struct {
		token  string
		userId int
		orgId  int
		email  string
		error  error
	}{
		{
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjk4LCJFbWFpbCI6Im1hdHRAZXhhbXBsZS5jb20iLCJPcmdhbml6YXRpb25JZCI6MjI0M30.Xb_scxUShDoaBRmKqxSzwrdPlZq8fzHNBy39fe6rrGI",
			98,
			2243,
			"matt@example.com",
			nil,
		},
		{
			"",
			-1,
			-1,
			"",
			errors.New("auth_missing_token"),
		},
		{
			"Be",
			-1,
			-1,
			"",
			errors.New("auth_malformed_token"),
		},
		{
			"BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjk4LCJFbWFpbCI6Im1hdHRAZXhhbXBsZS5jb20ifQ.zbpK1ZSeOTrsvSscE7KJqdhHfMgiOSHiu_2jOBsOLyA",
			-1,
			-1,
			"",
			errors.New("auth_malformed_token"),
		},
		{
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ2.eyJVc2VySWQiOjk4LCJFbWFpbCI6Im1hdHRAZXhhbXBsZS5jb20ifQ.zbpK1ZSeOTrsvSscE7KJqdhHfMgiOSHiu_2jOBsOLyA",
			-1,
			-1,
			"",
			errors.New("auth_malformed_token"),
		},
		{
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjk4LCJFbWFpbCI6Im1hdHRAZXhhbXBsZS5jb20ifQ.zbpK1ZSeOTrsvSscEdKJqdhHfMgiOSHiu_2jOBsOLyA",
			-1,
			-1,
			"",
			errors.New("auth_malformed_token"),
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		helper := setupBasicHelper(t, mockCtrl)

		// act
		userId, orgId, email, err := helper.ValidateToken(testCase.token)

		// assert
		assert.Equal(t, testCase.userId, userId)
		assert.Equal(t, testCase.orgId, orgId)
		assert.Equal(t, testCase.email, email)
		assert.Equal(t, testCase.error, err)
	}
}

func TestBasicHelper_ValidatePassword(t *testing.T) {
	testCases := []struct {
		password string
		error    error
	}{
		{
			"example1252352",
			nil,
		},
		{
			"abcdef",
			nil,
		},
		{
			"abcde",
			errors.New("validate_password_insufficient"),
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, testCase := range testCases {
		// arrange
		helper := setupBasicHelper(t, mockCtrl)

		// act
		error := helper.ValidatePassword(testCase.password)

		// assert
		assert.Equal(t, testCase.error, error)
	}
}
