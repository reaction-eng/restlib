package roles_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/reaction-eng/restlib/roles"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRepoMySql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT roleId FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("DELETE FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("INSERT INTO " + roles.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTable := mocks.NewMockPermissionTable(mockCtrl)

	// act
	repoMySql, err := roles.NewRepoMySql(db, mockTable)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func TestRepoPostgresSql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT roleId FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("DELETE FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("INSERT INTO " + roles.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTable := mocks.NewMockPermissionTable(mockCtrl)

	// act
	repoMySql, err := roles.NewRepoPostgresSql(db, mockTable)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func setupSqlMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("SELECT roleId FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("DELETE FROM " + roles.TableName + " WHERE userId = (.+) AND orgId = (.+)")
	mock.ExpectPrepare("INSERT INTO " + roles.TableName)

	return db, mock
}

func TestGetPermissions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment             string
		getUserRolesRows    *sqlmock.Rows
		getUserRolesError   error
		setupPermTable      func(table *mocks.MockPermissionTable)
		expectedPermissions *roles.Permissions
		expectedError       error
	}{
		{
			comment:           "working",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(10).AddRow(40),
			getUserRolesError: nil,
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().GetPermissions(10).Return([]string{"perm1"})
				mockPermTable.EXPECT().GetPermissions(40).Return([]string{"perm2", "perm3"})
			},
			expectedPermissions: &roles.Permissions{
				Permissions: []string{"perm1", "perm2", "perm3"},
			},
		},
		{
			comment:           "no roles",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}),
			getUserRolesError: nil,
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
			},
			expectedPermissions: &roles.Permissions{
				Permissions: []string{},
			},
		},
		{
			comment:           "db error",
			getUserRolesRows:  nil,
			getUserRolesError: errors.New("can't access db"),
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
			},
			expectedPermissions: nil,
			expectedError:       errors.New("can't access db"),
		},
		{
			comment:           "scan error",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow("alphea beta"),
			getUserRolesError: nil,
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
			},
			expectedPermissions: nil,
			expectedError:       errors.New("sql: Scan error on column index 0, name \"rowId\": converting driver.Value type string (\"alphea beta\") to a int: invalid syntax"),
		},
		{
			comment:           "row error",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(10).RowError(0, errors.New("db error")),
			getUserRolesError: nil,
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
			},
			expectedPermissions: nil,
			expectedError:       errors.New("db error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		mockTable := mocks.NewMockPermissionTable(mockCtrl)
		testCase.setupPermTable(mockTable)
		repo, err := roles.NewRepoPostgresSql(db, mockTable)

		user := mocks.NewMockUser(mockCtrl)
		user.EXPECT().Id().Times(1).Return(34)

		dbMock.ExpectQuery("SELECT roleId FROM  "+roles.TableName).
			WithArgs(34, 1000).
			WillReturnRows(testCase.getUserRolesRows).
			WillReturnError(testCase.getUserRolesError)

		// act
		permissions, err := repo.GetPermissions(user, 1000)

		// assert
		assert.Equal(t, testCase.expectedPermissions, permissions)
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestGetRoleIds(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment           string
		getUserRolesRows  *sqlmock.Rows
		getUserRolesError error
		expectedRoles     []int
		expectedError     error
	}{
		{
			comment:           "working",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(10).AddRow(40),
			getUserRolesError: nil,
			expectedRoles:     []int{10, 40},
			expectedError:     nil,
		},
		{
			comment:           "no roles",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}),
			getUserRolesError: nil,
			expectedRoles:     []int{},
			expectedError:     nil,
		},
		{
			comment:           "db error",
			getUserRolesRows:  nil,
			getUserRolesError: errors.New("can't access db"),
			expectedRoles:     nil,
			expectedError:     errors.New("can't access db"),
		},
		{
			comment:           "scan error",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow("alphea beta"),
			getUserRolesError: nil,
			expectedRoles:     nil,
			expectedError:     errors.New("sql: Scan error on column index 0, name \"rowId\": converting driver.Value type string (\"alphea beta\") to a int: invalid syntax"),
		},
		{
			comment:           "row error",
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(10).RowError(0, errors.New("db error")),
			getUserRolesError: nil,
			expectedRoles:     nil,
			expectedError:     errors.New("db error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		mockTable := mocks.NewMockPermissionTable(mockCtrl)
		repo, err := roles.NewRepoPostgresSql(db, mockTable)

		user := mocks.NewMockUser(mockCtrl)
		user.EXPECT().Id().Times(1).Return(34)

		dbMock.ExpectQuery("SELECT roleId FROM  "+roles.TableName).
			WithArgs(34, 1000).
			WillReturnRows(testCase.getUserRolesRows).
			WillReturnError(testCase.getUserRolesError)

		// act
		roles, err := repo.GetRoleIds(user, 1000)

		// assert
		assert.Equal(t, testCase.expectedRoles, roles)
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestSetRolesByRoleId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment           string
		inputRoles        []int
		getUserRolesRows  *sqlmock.Rows
		getUserRolesError error
		expectClear       bool
		clearUserRoles    error
		setupAddUserRole  func(sqlMock *sqlmock.Sqlmock)
		expectedError     error
	}{
		{
			comment:           "no update",
			inputRoles:        []int{3, 2, 5},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(2).AddRow(3).AddRow(5),
			getUserRolesError: nil,
			expectClear:       false,
			clearUserRoles:    nil,
			setupAddUserRole:  func(sqlMock *sqlmock.Sqlmock) {},
			expectedError:     nil,
		},
		{
			comment:           "update update",
			inputRoles:        []int{3, 2, 5, 6},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(2).AddRow(3).AddRow(5),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 3).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 5).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 6).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			comment:           "remove roles update",
			inputRoles:        []int{},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(2).AddRow(3).AddRow(5),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: nil,
		},
		{
			comment:           "no existing roles update",
			inputRoles:        []int{3, 2, 5, 6},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 3).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 5).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 6).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			comment:           "get user roles error",
			inputRoles:        []int{3, 2, 5, 6},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}),
			getUserRolesError: errors.New("get user role error"),
			clearUserRoles:    nil,
			expectClear:       false,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: errors.New("get user role error"),
		},
		{
			comment:           "can't clear",
			inputRoles:        []int{3, 2, 5, 6},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(4),
			getUserRolesError: nil,
			clearUserRoles:    errors.New("can't clear error"),
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: errors.New("can't clear error"),
		},
		{
			comment:           "can't add user rol error",
			inputRoles:        []int{3, 2, 5, 6},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(4),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 3).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(errors.New("can't add user error"))
			},
			expectedError: errors.New("can't add user error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		mockTable := mocks.NewMockPermissionTable(mockCtrl)
		repo, err := roles.NewRepoPostgresSql(db, mockTable)

		user := mocks.NewMockUser(mockCtrl)
		user.EXPECT().Id().Times(2).Return(34)

		dbMock.ExpectQuery("SELECT roleId FROM  "+roles.TableName).
			WithArgs(34, 1000).
			WillReturnRows(testCase.getUserRolesRows).
			WillReturnError(testCase.getUserRolesError)

		if testCase.expectClear {
			dbMock.ExpectExec("DELETE  FROM  "+roles.TableName).
				WithArgs(34, 1000).
				WillReturnResult(sqlmock.NewResult(0, 2)).
				WillReturnError(testCase.clearUserRoles)
		}

		testCase.setupAddUserRole(&dbMock)

		// act
		err = repo.SetRolesByRoleId(user, 1000, testCase.inputRoles)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestSetRolesByRoleName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		comment           string
		setupPermTable    func(table *mocks.MockPermissionTable)
		inputRoles        []string
		getUserRolesRows  *sqlmock.Rows
		getUserRolesError error
		expectClear       bool
		clearUserRoles    error
		setupAddUserRole  func(sqlMock *sqlmock.Sqlmock)
		expectedError     error
	}{
		{
			comment:    "no update",
			inputRoles: []string{"role4", "role 1", "role 2"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role4").Return(4, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 1").Return(1, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(1).AddRow(2).AddRow(4),
			getUserRolesError: nil,
			expectClear:       false,
			clearUserRoles:    nil,
			setupAddUserRole:  func(sqlMock *sqlmock.Sqlmock) {},
			expectedError:     nil,
		},
		{
			comment:    "update update",
			inputRoles: []string{"role4", "role 1", "role 2"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role4").Return(4, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 1").Return(1, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(2).AddRow(3).AddRow(5),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 1).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 4).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			comment:    "remove roles update",
			inputRoles: []string{},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(2).AddRow(3).AddRow(5),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: nil,
		},
		{
			comment:    "no existing roles update",
			inputRoles: []string{"role 3", "role 2", "role 5", "role 6"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role 3").Return(3, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 5").Return(5, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 6").Return(6, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(5),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 3).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 2).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 5).WillReturnResult(sqlmock.NewResult(0, 1))
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 6).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			comment:    "get user roles error",
			inputRoles: []string{"role 3", "role 2", "role 5", "role 6"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role 3").Return(3, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 5").Return(5, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 6").Return(6, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}),
			getUserRolesError: errors.New("get user role error"),
			clearUserRoles:    nil,
			expectClear:       false,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: errors.New("get user role error"),
		},
		{
			comment:    "can't clear",
			inputRoles: []string{"role 3", "role 2", "role 5", "role 6"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role 3").Return(3, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 5").Return(5, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 6").Return(6, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(4),
			getUserRolesError: nil,
			clearUserRoles:    errors.New("can't clear error"),
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
			},
			expectedError: errors.New("can't clear error"),
		},
		{
			comment:    "can't add user rol error",
			inputRoles: []string{"role 3", "role 2", "role 5", "role 6"},
			setupPermTable: func(mockPermTable *mocks.MockPermissionTable) {
				mockPermTable.EXPECT().LookUpRoleId("role 3").Return(3, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 2").Return(2, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 5").Return(5, nil)
				mockPermTable.EXPECT().LookUpRoleId("role 6").Return(6, nil)
			},
			getUserRolesRows:  sqlmock.NewRows([]string{"rowId"}).AddRow(4),
			getUserRolesError: nil,
			clearUserRoles:    nil,
			expectClear:       true,
			setupAddUserRole: func(sqlMock *sqlmock.Sqlmock) {
				(*sqlMock).ExpectExec("INSERT INTO "+roles.TableName).WithArgs(34, 1000, 3).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(errors.New("can't add user error"))
			},
			expectedError: errors.New("can't add user error"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t)

		mockTable := mocks.NewMockPermissionTable(mockCtrl)
		repo, err := roles.NewRepoPostgresSql(db, mockTable)

		user := mocks.NewMockUser(mockCtrl)
		user.EXPECT().Id().Times(2).Return(34)

		testCase.setupPermTable(mockTable)

		dbMock.ExpectQuery("SELECT roleId FROM  "+roles.TableName).
			WithArgs(34, 1000).
			WillReturnRows(testCase.getUserRolesRows).
			WillReturnError(testCase.getUserRolesError)

		if testCase.expectClear {
			dbMock.ExpectExec("DELETE  FROM  "+roles.TableName).
				WithArgs(34, 1000).
				WillReturnResult(sqlmock.NewResult(0, 2)).
				WillReturnError(testCase.clearUserRoles)
		}

		testCase.setupAddUserRole(&dbMock)

		// act
		err = repo.SetRolesByName(user, 1000, testCase.inputRoles)

		// assert
		assert.Equal(t, testCase.expectedError, err)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}
