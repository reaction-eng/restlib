package preferences_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/reaction-eng/restlib/mocks"

	"github.com/reaction-eng/restlib/preferences"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewRepoMySql(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT settings FROM " + preferences.TableName)
	mock.ExpectPrepare("INSERT INTO " + preferences.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	optionsGroup := &preferences.OptionGroup{}

	// act
	repoMySql, err := preferences.NewRepoMySql(db, optionsGroup)

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

	mock.ExpectPrepare("SELECT settings FROM " + preferences.TableName)
	mock.ExpectPrepare("INSERT INTO " + preferences.TableName)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	optionsGroup := &preferences.OptionGroup{}

	// act
	repoMySql, err := preferences.NewRepoPostgresSql(db, optionsGroup)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, repoMySql)
}

func setupSqlMock(t *testing.T, mockCtrl *gomock.Controller) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectPrepare("SELECT settings FROM " + preferences.TableName)
	mock.ExpectPrepare("INSERT INTO " + preferences.TableName)

	return db, mock
}

func TestRepoSql_GetPreferences(t *testing.T) {
	optionsGroup := &preferences.OptionGroup{
		Options: []preferences.Option{
			{
				Id:           "id1",
				DefaultValue: "33",
			},
			{
				Id: "id2",
			},
		},
	}

	testCases := []struct {
		comment             string
		settingString       string
		queryError          error
		expectedPreferences *preferences.Preferences
		expectedError       error
	}{
		{
			comment:       "working settings and no errors",
			settingString: `{"settings":{"id1":"true","id2":"blue"}}`,
			queryError:    nil,
			expectedPreferences: &preferences.Preferences{
				Options: optionsGroup,
				Settings: &preferences.SettingGroup{
					Settings: map[string]string{"id1": "true", "id2": "blue"},
					SubGroup: map[string]*preferences.SettingGroup{},
				},
			},
			expectedError: nil,
		}, {
			comment:       "no rows, returns default values",
			settingString: ``,
			queryError:    sql.ErrNoRows,
			expectedPreferences: &preferences.Preferences{
				Options: optionsGroup,
				Settings: &preferences.SettingGroup{
					Settings: map[string]string{"id1": "33", "id2": ""},
					SubGroup: map[string]*preferences.SettingGroup{},
				},
			},
			expectedError: nil,
		},
		{
			comment:             "db error, pass along",
			settingString:       ``,
			queryError:          errors.New("new db error"),
			expectedPreferences: nil,
			expectedError:       errors.New("new db error"),
		},
		{
			comment:       "should repair missing parameters",
			settingString: `{"settings":{"id2":"blue"}}`,
			queryError:    nil,
			expectedPreferences: &preferences.Preferences{
				Options: optionsGroup,
				Settings: &preferences.SettingGroup{
					Settings: map[string]string{"id1": "33", "id2": "blue"},
					SubGroup: map[string]*preferences.SettingGroup{},
				},
			},
			expectedError: nil,
		},
		{
			comment:             "should throw error for bad json",
			settingString:       `{{`,
			queryError:          nil,
			expectedPreferences: nil,
			expectedError:       errors.New(`sql: Scan error on column index 0, name "setting": invalid character '{' looking for beginning of object key string`),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t, mockCtrl)

		mockUser := mocks.NewMockUser(mockCtrl)
		mockUser.EXPECT().Id().Return(3).MinTimes(1)

		var rows *sqlmock.Rows
		if len(testCase.settingString) > 0 {
			rows = sqlmock.NewRows([]string{"setting"}).
				AddRow(testCase.settingString)
		}

		dbMock.ExpectQuery("SELECT settings FROM " + preferences.TableName).
			WithArgs(mockUser.Id()).
			WillReturnRows(rows).
			WillReturnError(testCase.queryError)

		repo, err := preferences.NewRepoPostgresSql(db, optionsGroup)
		assert.Nil(t, err)

		// act
		preference, err := repo.GetPreferences(mockUser)
		assert.Equal(t, testCase.expectedError, err)

		// assert
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, testCase.expectedPreferences, preference)

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}

func TestRepoSql_SetPreferences(t *testing.T) {
	optionsGroup := &preferences.OptionGroup{}

	testCases := []struct {
		comment             string
		expectedPreferences *preferences.Preferences
		execError           error
		expectedError       error
	}{
		{
			comment: "working settings and no errors",
			expectedPreferences: &preferences.Preferences{
				Options: optionsGroup,
				Settings: &preferences.SettingGroup{
					Settings: map[string]string{"id1": "true", "id2": "blue"},
					SubGroup: map[string]*preferences.SettingGroup{},
				},
			},
			execError:     nil,
			expectedError: nil,
		},
		{
			comment: "with error",
			expectedPreferences: &preferences.Preferences{
				Options: optionsGroup,
				Settings: &preferences.SettingGroup{
					Settings: map[string]string{"id1": "true", "id2": "blue"},
					SubGroup: map[string]*preferences.SettingGroup{},
				},
			},
			execError:     errors.New("error from db set"),
			expectedError: errors.New("error from db set"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		mockCtrl := gomock.NewController(t)

		db, dbMock := setupSqlMock(t, mockCtrl)

		mockUser := mocks.NewMockUser(mockCtrl)
		mockUser.EXPECT().Id().Return(3).MinTimes(1)

		dbMock.ExpectExec("INSERT INTO "+preferences.TableName).
			WithArgs(mockUser.Id(), testCase.expectedPreferences.Settings).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(testCase.execError)

		repo, err := preferences.NewRepoPostgresSql(db, optionsGroup)
		assert.Nil(t, err)

		// act
		preference, err := repo.SetPreferences(mockUser, testCase.expectedPreferences.Settings)

		// assert
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, testCase.expectedPreferences, preference)
		assert.Equal(t, testCase.expectedError, err)

		// cleanup
		db.Close()
		mockCtrl.Finish()
	}
}
