// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/reaction-eng/restlib/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/users"
	"github.com/stretchr/testify/assert"
)

func TestNewRepoMySql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("SELECT id, activation FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT orgId FROM " + users.UserOrgTableName)
	mock.ExpectPrepare("INSERT INTO " + users.UserOrgTableName)
	mock.ExpectPrepare("DELETE FROM " + users.UserOrgTableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// act
	repoMySql, err := users.NewRepoMySql(db)

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

	mock.ExpectPrepare("INSERT INTO " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("SELECT id, activation FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT orgId FROM " + users.UserOrgTableName)
	mock.ExpectPrepare("INSERT INTO " + users.UserOrgTableName)
	mock.ExpectPrepare("DELETE FROM " + users.UserOrgTableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// act
	repoMySql, err := users.NewRepoPostgresSql(db)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func setupSqlMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("INSERT INTO " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT \\* FROM " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("UPDATE " + users.UserTableName)
	mock.ExpectPrepare("SELECT id, activation FROM " + users.UserTableName)
	mock.ExpectPrepare("SELECT orgId FROM " + users.UserOrgTableName)
	mock.ExpectPrepare("INSERT INTO " + users.UserOrgTableName)
	mock.ExpectPrepare("DELETE FROM " + users.UserOrgTableName)

	return db, mock
}

func TestResetRepoSql_GetUserByEmail(t *testing.T) {
	testCases := []struct {
		comment           string
		emailInput        string
		expectedError     error
		query             func(*sqlmock.ExpectedQuery)
		queryOrgs         func(*sqlmock.ExpectedQuery)
		userId            int
		userOrgs          []int
		userEmail         string
		userToken         string
		userActivated     bool
		userPasswordLogin bool
	}{
		{
			comment:       "working",
			emailInput:    "user@example.com",
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: true,
		},
		{
			comment: "crazy email input",
			emailInput: "  	usEr@example.Com		",
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: true,
		},
		{
			comment:       "no rows",
			emailInput:    "user@example.com",
			expectedError: users.UserNotFound,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").WillReturnError(sql.ErrNoRows)
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
		},
		{
			comment:       "nil date should be non activated",
			emailInput:    "user@example.com",
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", nil))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     false,
			userPasswordLogin: true,
		},
		{
			comment:       "empty password should prevent user from logging in",
			emailInput:    "user@example.com",
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: false,
		},
		{
			comment:       "empty password and nil date should trigger bools",
			emailInput:    "user@example.com",
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "", nil))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     false,
			userPasswordLogin: false,
		},
		{
			comment:       "get org error should return nil",
			emailInput:    "user@example.com",
			expectedError: errors.New("could not get org info"),
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.com").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23)).
					WillReturnError(errors.New("could not get org info"))
			},
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		repo, err := users.NewRepoPostgresSql(db)
		testCase.query(dbMock.ExpectQuery("SELECT \\* FROM users"))
		testCase.queryOrgs(dbMock.ExpectQuery("SELECT orgId FROM userOrganizations"))
		assert.Nil(t, err)

		// act
		user, err := repo.GetUserByEmail(testCase.emailInput)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if err == nil {
			assert.Equal(t, testCase.userId, user.Id())
			assert.Equal(t, testCase.userOrgs, user.Organizations())
			assert.Equal(t, testCase.userEmail, user.Email())
			assert.Equal(t, testCase.userActivated, user.Activated())
			assert.Equal(t, testCase.userPasswordLogin, user.PasswordLogin())
			assert.Equal(t, testCase.userToken, user.Token())
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_GetUser(t *testing.T) {
	testCases := []struct {
		comment           string
		idInput           int
		expectedError     error
		query             func(*sqlmock.ExpectedQuery)
		queryOrgs         func(*sqlmock.ExpectedQuery)
		userId            int
		userOrgs          []int
		userEmail         string
		userToken         string
		userActivated     bool
		userPasswordLogin bool
	}{
		{
			comment:       "working",
			idInput:       43,
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: true,
		},
		{
			comment:       "other db error",
			idInput:       43,
			expectedError: errors.New("other db error"),
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now())).
					WillReturnError(errors.New("other db error"))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: true,
		},
		{
			comment:       "no rows",
			idInput:       43,
			expectedError: users.UserNotFound,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).WillReturnError(sql.ErrNoRows)
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
		},
		{
			comment:       "nil date should be non activated",
			idInput:       43,
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", nil))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23).AddRow(32))
			},
			userId:            43,
			userOrgs:          []int{54, 23, 32},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     false,
			userPasswordLogin: true,
		},
		{
			comment:       "empty password should prevent user from logging in",
			idInput:       43,
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: false,
		},
		{
			comment:       "empty password and nil date should trigger bools",
			idInput:       43,
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "", nil))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}).
							AddRow(54).AddRow(23))
			},
			userId:            43,
			userOrgs:          []int{54, 23},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     false,
			userPasswordLogin: false,
		},
		{
			comment:       "no orgs should work",
			idInput:       43,
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}))
			},
			userId:            43,
			userOrgs:          []int{},
			userEmail:         "user@example.com",
			userToken:         "",
			userActivated:     true,
			userPasswordLogin: true,
		},
		{
			comment:       "org db error",
			idInput:       43,
			expectedError: errors.New("org db error"),
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.com", "password", time.Now()))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"})).
					WillReturnError(errors.New("org db error"))
			},
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.query(dbMock.ExpectQuery("SELECT \\* FROM users"))
		testCase.queryOrgs(dbMock.ExpectQuery("SELECT orgId FROM userOrganizations"))

		// act
		user, err := repo.GetUser(testCase.idInput)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if err == nil {
			assert.Equal(t, testCase.userId, user.Id())
			assert.Equal(t, testCase.userOrgs, user.Organizations())
			assert.Equal(t, testCase.userEmail, user.Email())
			assert.Equal(t, testCase.userActivated, user.Activated())
			assert.Equal(t, testCase.userPasswordLogin, user.PasswordLogin())
			assert.Equal(t, testCase.userToken, user.Token())
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_ListUsers(t *testing.T) {
	testCases := []struct {
		comment       string
		onlyActive    bool
		orgList       []int //TODO: add org test
		expectedUsers []int
		expectedError error
		query         func(query *sqlmock.ExpectedQuery)
	}{
		{
			comment:       "working",
			onlyActive:    false,
			orgList:       nil,
			expectedUsers: []int{32, 52, 2432, 23},
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "activationDate"}).
							AddRow(32, time.Now()).
							AddRow(52, time.Now()).
							AddRow(2432, time.Now()).
							AddRow(23, time.Now()))
			},
		},
		{
			comment:       "db broken",
			onlyActive:    false,
			orgList:       nil,
			expectedUsers: nil,
			expectedError: errors.New("db broken"),
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnError(errors.New("db broken"))
			},
		},
		{
			comment:       "row errors out",
			onlyActive:    false,
			orgList:       nil,
			expectedUsers: []int{32, 52},
			expectedError: errors.New("row error"),
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "activationDate"}).
							AddRow(32, time.Now()).
							AddRow(52, time.Now()).RowError(2, errors.New("row error")).
							AddRow(2432, time.Now()).
							AddRow(23, time.Now()))
			},
		},
		{
			comment:       "empty list",
			onlyActive:    false,
			orgList:       nil,
			expectedUsers: []int{},
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "activationDate"}))
			},
		},
		{
			comment:       "with both activated and non active",
			onlyActive:    false,
			orgList:       nil,
			expectedUsers: []int{32, 52, 2432, 23},
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "activationDate"}).
							AddRow(32, time.Now()).
							AddRow(52, nil).
							AddRow(2432, nil).
							AddRow(23, time.Now()))
			},
		},
		{
			comment:       "only active",
			onlyActive:    true,
			orgList:       nil,
			expectedUsers: []int{32, 23},
			expectedError: nil,
			query: func(query *sqlmock.ExpectedQuery) {
				query.
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "activationDate"}).
							AddRow(32, time.Now()).
							AddRow(52, nil).
							AddRow(2432, nil).
							AddRow(23, time.Now()))
			},
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.query(dbMock.ExpectQuery("SELECT id, activation FROM users"))

		// act
		userList, err := repo.ListUsers(testCase.onlyActive, testCase.orgList)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUsers, userList)

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestResetRepoSql_AddUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment                string
		expectedError          error
		inputUser              func() users.User
		checkUserStatementExec func(exec *sqlmock.ExpectedQuery)
		checkUserQueryOrgs     func(*sqlmock.ExpectedQuery)
		addUserStatementExec   func(exec *sqlmock.ExpectedExec)
		getUserByEmailQuery    func(query *sqlmock.ExpectedQuery)
		queryOrgs              func(*sqlmock.ExpectedQuery)
		userId                 int
		userEmail              string
		userToken              string
		userActivated          bool
		userPasswordLogin      bool
	}{
		{
			comment: "working",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(3)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			expectedError: nil,
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnError(sql.ErrNoRows)
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {
				exec.
					WithArgs("user@example.info", "hashed password").
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.info", "password", nil))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}))
			},
			userId:            43,
			userEmail:         "user@example.info",
			userToken:         "",
			userActivated:     false,
			userPasswordLogin: true,
		},
		{
			comment: "user already exists",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(1)
				user.EXPECT().Password().Return("hashed password").Times(0)
				return user
			},
			expectedError: errors.New("user_email_in_user"),
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.info", "password", nil))
			},
			checkUserQueryOrgs: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs(43).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"orgId"}))
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {

			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {

			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
			},
		},
		{
			comment: "add user statement error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(2)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			expectedError: errors.New("db error"),
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.info", "password", nil)).
					WillReturnError(sql.ErrNoRows)
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {
				exec.
					WithArgs("user@example.info", "hashed password").
					WillReturnResult(sqlmock.NewResult(3, 3)).
					WillReturnError(errors.New("db error"))
			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {

			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {
			},
		},
		{
			comment: "get user by email error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(3)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			expectedError: errors.New("db error"),
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnError(sql.ErrNoRows)
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {
				exec.
					WithArgs("user@example.info", "hashed password").
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.info", "password", nil)).
					WillReturnError(errors.New("db error"))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {

			},
		},
		{
			comment: "get user by email error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(3)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			expectedError: errors.New("db error"),
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnError(sql.ErrNoRows)
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {
				exec.
					WithArgs("user@example.info", "hashed password").
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "email", "password", "activationDate"}).
							AddRow(43, "user@example.info", "password", nil)).
					WillReturnError(errors.New("db error"))
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {

			},
		},
		{
			comment: "get user by email error in check",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Email().Return("user@example.info").Times(1)
				user.EXPECT().Password().Return("hashed password").Times(0)
				return user
			},
			expectedError: errors.New("db error"),
			checkUserStatementExec: func(query *sqlmock.ExpectedQuery) {
				query.
					WithArgs("user@example.info").
					WillReturnError(errors.New("db error"))
			},
			addUserStatementExec: func(exec *sqlmock.ExpectedExec) {
			},
			getUserByEmailQuery: func(query *sqlmock.ExpectedQuery) {
			},
			queryOrgs: func(query *sqlmock.ExpectedQuery) {

			},
		},
	}

	for _, testCase := range testCases {
		// arrange
		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.checkUserStatementExec(dbMock.ExpectQuery("SELECT \\* FROM users "))
		if testCase.checkUserQueryOrgs != nil {
			testCase.checkUserQueryOrgs(dbMock.ExpectQuery("SELECT orgId FROM userOrganizations"))
		}
		testCase.addUserStatementExec(dbMock.ExpectExec("INSERT INTO users "))
		testCase.getUserByEmailQuery(dbMock.ExpectQuery("SELECT \\* FROM users "))
		testCase.queryOrgs(dbMock.ExpectQuery("SELECT orgId FROM userOrganizations"))

		// act
		user, err := repo.AddUser(testCase.inputUser())

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if err == nil {
			assert.Equal(t, testCase.userId, user.Id())
			assert.Equal(t, testCase.userEmail, user.Email())
			assert.Equal(t, testCase.userActivated, user.Activated())
			assert.Equal(t, testCase.userPasswordLogin, user.PasswordLogin())
			assert.Equal(t, testCase.userToken, user.Token())
		}

		// cleanup
		db.Close()
	}
}

func TestResetRepoSql_UpdateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment        string
		inputUser      func() users.User
		updateUserExec func(exec *sqlmock.ExpectedExec)
		expectedError  error
	}{
		{
			comment: "working",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				user.EXPECT().Email().Return("user@example.info").Times(1)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			updateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs("user@example.info", "hashed password", 34).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			expectedError: nil,
		},
		{
			comment: "update user error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				user.EXPECT().Email().Return("user@example.info").Times(1)
				user.EXPECT().Password().Return("hashed password").Times(1)
				return user
			},
			updateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs("user@example.info", "hashed password", 34).
					WillReturnResult(sqlmock.NewResult(3, 3)).
					WillReturnError(errors.New("db error"))

			},
			expectedError: errors.New("db error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.updateUserExec(dbMock.ExpectExec("UPDATE user "))

		userInput := testCase.inputUser()

		// act
		user, err := repo.UpdateUser(userInput)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, userInput, user)

		// cleanup
		db.Close()
	}
}

func TestResetRepoSql_ActivateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment          string
		inputUser        func() users.User
		activateUserExec func(exec *sqlmock.ExpectedExec)
		expectedError    error
	}{
		{
			comment: "working",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(sqlmock.AnyArg(), 34).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			expectedError: nil,
		},
		{
			comment: "activate with error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(sqlmock.AnyArg(), 34).
					WillReturnResult(sqlmock.NewResult(3, 3)).
					WillReturnError(errors.New("db error"))

			},
			expectedError: errors.New("db error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.activateUserExec(dbMock.ExpectExec("UPDATE user "))

		userInput := testCase.inputUser()

		// act
		err := repo.ActivateUser(userInput)

		// assert
		assert.Equal(t, testCase.expectedError, err)

		// cleanup
		db.Close()
	}
}

func TestResetRepoSql_AddUserToOrganization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment          string
		inputUser        func() users.User
		orgId            int
		activateUserExec func(exec *sqlmock.ExpectedExec)
		expectedError    error
	}{
		{
			comment: "working",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			orgId:         3234,
			expectedError: nil,
		},
		{
			comment: "add with error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(3, 3)).
					WillReturnError(errors.New("db error"))

			},
			orgId:         3234,
			expectedError: errors.New("db error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.activateUserExec(dbMock.ExpectExec("INSERT INTO userOrganizations "))

		userInput := testCase.inputUser()

		// act
		err := repo.AddUserToOrganization(userInput, testCase.orgId)

		// assert
		assert.Equal(t, testCase.expectedError, err)

		// cleanup
		db.Close()
	}
}

func TestResetRepoSql_RemoveUserFromOrganization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment          string
		inputUser        func() users.User
		orgId            int
		activateUserExec func(exec *sqlmock.ExpectedExec)
		expectedError    error
	}{
		{
			comment: "working",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			orgId:         3234,
			expectedError: nil,
		},
		{
			comment: "add with error",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234).
					WillReturnResult(sqlmock.NewResult(3, 3)).
					WillReturnError(errors.New("db error"))

			},
			orgId:         3234,
			expectedError: errors.New("db error"),
		},
		{
			comment: "no rows errors",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234).
					WillReturnResult(sqlmock.NewResult(3, 0))
			},
			orgId:         3234,
			expectedError: errors.New("no_organizations_removed"),
		}, {
			comment: "errors from results",
			inputUser: func() users.User {
				user := mocks.NewMockUser(mockCtrl)
				user.EXPECT().Id().Return(34).Times(1)
				return user
			},
			activateUserExec: func(query *sqlmock.ExpectedExec) {
				query.
					WithArgs(34, 3234).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("errors from results")))
			},
			orgId:         3234,
			expectedError: errors.New("errors from results"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		db, dbMock := setupSqlMock(t)

		repo, _ := users.NewRepoPostgresSql(db)
		testCase.activateUserExec(dbMock.ExpectExec("DELETE FROM userOrganizations "))

		userInput := testCase.inputUser()

		// act
		err := repo.RemoveUserFromOrganization(userInput, testCase.orgId)

		// assert
		assert.Equal(t, testCase.expectedError, err)

		// cleanup
		db.Close()
	}
}
