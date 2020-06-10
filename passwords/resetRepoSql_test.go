// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/reaction-eng/restlib/email"

	"github.com/reaction-eng/restlib/passwords"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewRepoMySql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO " + passwords.TableName)
	mock.ExpectPrepare("SELECT (.+) FROM " + passwords.TableName)
	mock.ExpectPrepare("delete FROM " + passwords.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEmailer := mocks.NewMockEmailer(mockCtrl)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStruct("password_reset", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("user_activation", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("one_time_password", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "one time email subject"
		as.Template = "one time email template"
	})
	mockConfiguration.EXPECT().GetFloat("tokenLifeSpan").Times(1).Return(float64(24), nil)

	// act
	repoMySql, err := passwords.NewRepoMySql(db, mockEmailer, mockConfiguration)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func TestNewRepoPostgresSql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO " + passwords.TableName)
	mock.ExpectPrepare("SELECT (.+) FROM " + passwords.TableName)
	mock.ExpectPrepare("delete FROM " + passwords.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEmailer := mocks.NewMockEmailer(mockCtrl)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStruct("password_reset", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("user_activation", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("one_time_password", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "one time email subject"
		as.Template = "one time email template"
	})
	mockConfiguration.EXPECT().GetFloat("tokenLifeSpan").Times(1).Return(float64(24), nil)

	// act
	repoMySql, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func setupSqlMock(t *testing.T, mockCtrl *gomock.Controller, tableName string) (*sql.DB, sqlmock.Sqlmock, *mocks.MockEmailer, *mocks.MockConfiguration) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("INSERT INTO " + tableName)
	mock.ExpectPrepare("SELECT (.+) FROM " + tableName)
	mock.ExpectPrepare("delete FROM " + tableName)

	mockEmailer := mocks.NewMockEmailer(mockCtrl)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStruct("password_reset", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "password_reset_subject"
		as.Template = "password_reset_template"
	})
	mockConfiguration.EXPECT().GetStruct("user_activation", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "user_activation_subject"
		as.Template = "user_activation_template"
	})
	mockConfiguration.EXPECT().GetStruct("one_time_password", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "one time email subject"
		as.Template = "one time email template"
	})
	mockConfiguration.EXPECT().GetFloat("tokenLifeSpan").Times(1).Return(float64(24), nil)

	return db, mock, mockEmailer, mockConfiguration
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestResetRepoSql_IssueResetRequest(t *testing.T) {
	testCases := []struct {
		token        string
		userId       int
		emailAddress string
		execError    error
		emailError   error
	}{
		{
			"exampleToken123",
			100,
			"test@example.com",
			nil,
			nil,
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			errors.New("exec error"),
			nil,
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			errors.New("exec error"),
			errors.New("email error"),
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			nil,
			errors.New("email error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		emailHeader := email.HeaderInfo{
			Subject: "password_reset_subject",
			To:      []string{testCase.emailAddress},
		}

		resetInfo := passwords.PasswordResetInfo{
			Token: testCase.token,
			Email: testCase.emailAddress,
		}

		if testCase.execError == nil { // this should only be called if the test case is nil
			mockEmailer.EXPECT().SendTemplateFile(&emailHeader, "password_reset_template", resetInfo, nil).Times(1).Return(testCase.emailError)
		}

		dbMock.ExpectExec("INSERT INTO "+passwords.TableName).
			WithArgs(testCase.userId, testCase.emailAddress, testCase.token, AnyTime{}, 2).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(testCase.execError)

		// act
		err = repo.IssueResetRequest(testCase.token, testCase.userId, testCase.emailAddress)

		// assert
		if testCase.execError != nil {
			assert.Equal(t, testCase.execError, err)
		} else if testCase.emailError != nil {
			assert.Equal(t, testCase.emailError, err)
		} else {
			assert.Nil(t, err)
		}
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_IssueActivationRequest(t *testing.T) {
	testCases := []struct {
		token        string
		userId       int
		emailAddress string
		execError    error
		emailError   error
	}{
		{
			"exampleToken123",
			100,
			"test@example.com",
			nil,
			nil,
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			errors.New("exec error"),
			nil,
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			errors.New("exec error"),
			errors.New("email error"),
		},
		{
			"exampleToken123",
			100,
			"test@example.com",
			nil,
			errors.New("email error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		emailHeader := email.HeaderInfo{
			Subject: "user_activation_subject",
			To:      []string{testCase.emailAddress},
		}

		resetInfo := passwords.PasswordResetInfo{
			Token: testCase.token,
			Email: testCase.emailAddress,
		}

		if testCase.execError == nil { // this should only be called if the test case is nil
			mockEmailer.EXPECT().SendTemplateFile(&emailHeader, "user_activation_template", resetInfo, nil).Times(1).Return(testCase.emailError)
		}

		dbMock.ExpectExec("INSERT INTO "+passwords.TableName).
			WithArgs(testCase.userId, testCase.emailAddress, testCase.token, AnyTime{}, 1). // one for activation
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(testCase.execError)

		// act
		err = repo.IssueActivationRequest(testCase.token, testCase.userId, testCase.emailAddress)

		// assert
		if testCase.execError != nil {
			assert.Equal(t, testCase.execError, err)
		} else if testCase.emailError != nil {
			assert.Equal(t, testCase.emailError, err)
		} else {
			assert.Nil(t, err)
		}
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_CheckForResetToken(t *testing.T) {
	testCases := []struct {
		userId        int
		token         string
		userIdDb      int
		tokenDb       string
		issueTime     time.Time
		queryError    error
		rowId         int
		expectedRowId int
		expectedError error
	}{
		{
			100,
			"example token",
			100,
			"example token",
			time.Now().Add(-25 * time.Hour),
			nil,
			1023,
			0,
			errors.New("token_expired"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			nil,
			1023,
			1023,
			nil,
		},
		{
			100,
			"example token",
			100,
			"other token token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("password_change_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"other example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("password_change_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("password_change_forbidden"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			errors.New("queryError"),
			1023,
			-1,
			errors.New("password_change_forbidden"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		rows := sqlmock.NewRows([]string{"id", "userId", "email", "token", "issued", "type"}).
			AddRow(testCase.rowId, testCase.userIdDb, "email", testCase.tokenDb, testCase.issueTime, 1)

		dbMock.ExpectQuery("SELECT (.+) FROM " + passwords.TableName).
			WillReturnRows(rows).
			WillReturnError(testCase.queryError)

		// act
		returnedRowId, err := repo.CheckForResetToken(testCase.userId, testCase.token)

		// assert
		assert.Equal(t, testCase.expectedRowId, returnedRowId)
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_CheckForActivationToken(t *testing.T) {
	testCases := []struct {
		userId        int
		token         string
		userIdDb      int
		tokenDb       string
		issueTime     time.Time
		queryError    error
		rowId         int
		expectedRowId int
		expectedError error
	}{
		{
			100,
			"example token",
			100,
			"example token",
			time.Now().Add(-25 * time.Hour),
			nil,
			1023,
			0,
			errors.New("token_expired"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			nil,
			1023,
			1023,
			nil,
		},
		{
			100,
			"example token",
			100,
			"other token token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("activation_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"other example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("activation_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("activation_forbidden"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			errors.New("queryError"),
			1023,
			-1,
			errors.New("activation_forbidden"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		rows := sqlmock.NewRows([]string{"id", "userId", "email", "token", "issued", "type"}).
			AddRow(testCase.rowId, testCase.userIdDb, "email", testCase.tokenDb, testCase.issueTime, 2)

		dbMock.ExpectQuery("SELECT (.+) FROM " + passwords.TableName).
			WillReturnRows(rows).
			WillReturnError(testCase.queryError)

		// act
		returnedRowId, err := repo.CheckForActivationToken(testCase.userId, testCase.token)

		// assert
		assert.Equal(t, testCase.expectedRowId, returnedRowId)
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_CheckForOneTimePasswordToken(t *testing.T) {
	testCases := []struct {
		userId        int
		token         string
		userIdDb      int
		tokenDb       string
		issueTime     time.Time
		queryError    error
		rowId         int
		expectedRowId int
		expectedError error
	}{
		{
			100,
			"example token",
			100,
			"example token",
			time.Now().Add(-25 * time.Hour),
			nil,
			1023,
			0,
			errors.New("token_expired"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			nil,
			1023,
			1023,
			nil,
		},
		{
			100,
			"example token",
			100,
			"other token token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("oneTimePassword_login_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"other example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("oneTimePassword_login_forbidden"),
		},
		{
			100,
			"example token",
			102,
			"example token",
			time.Now(),
			nil,
			1023,
			-1,
			errors.New("oneTimePassword_login_forbidden"),
		},
		{
			100,
			"example token",
			100,
			"example token",
			time.Now(),
			errors.New("queryError"),
			1023,
			-1,
			errors.New("oneTimePassword_login_forbidden"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, err := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		rows := sqlmock.NewRows([]string{"id", "userId", "email", "token", "issued", "type"}).
			AddRow(testCase.rowId, testCase.userIdDb, "email", testCase.tokenDb, testCase.issueTime, 3)

		dbMock.ExpectQuery("SELECT (.+) FROM " + passwords.TableName).
			WillReturnRows(rows).
			WillReturnError(testCase.queryError)

		// act
		returnedRowId, err := repo.CheckForOneTimePasswordToken(testCase.userId, testCase.token)

		// assert
		assert.Equal(t, testCase.expectedRowId, returnedRowId)
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_UseToken(t *testing.T) {
	testCases := []struct {
		tokenId int
		error   error
	}{
		{
			100,
			nil,
		},
		{
			100,
			errors.New("exampleError"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock, mockEmailer, mockConfiguration := setupSqlMock(t, mockCtrl, passwords.TableName)

		repo, _ := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

		dbMock.ExpectExec("delete FROM " + passwords.TableName).
			WithArgs(testCase.tokenId). // one for activation
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(testCase.error)

		// act
		err := repo.UseToken(testCase.tokenId)

		// assert
		assert.Equal(t, testCase.error, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_CleanUp(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMock.ExpectPrepare("INSERT INTO " + passwords.TableName).WillBeClosed()
	dbMock.ExpectPrepare("SELECT (.+) FROM " + passwords.TableName).WillBeClosed()
	dbMock.ExpectPrepare("delete FROM " + passwords.TableName).WillBeClosed()

	mockEmailer := mocks.NewMockEmailer(mockCtrl)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetStruct("password_reset", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("user_activation", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "test email subject"
		as.Template = "test email template"
	})
	mockConfiguration.EXPECT().GetStruct("one_time_password", gomock.Any()).Times(1).Do(func(name string, s interface{}) {
		as, _ := s.(*passwords.PasswordResetConfig)
		as.Subject = "one time email subject"
		as.Template = "one time email template"
	})
	mockConfiguration.EXPECT().GetFloat("tokenLifeSpan").Times(1).Return(float64(24), nil)

	repo, _ := passwords.NewRepoPostgresSql(db, mockEmailer, mockConfiguration)

	// act
	repo.CleanUp()

	// assert
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
