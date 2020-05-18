package roles_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/reaction-eng/restlib/roles"
	"github.com/stretchr/testify/assert"
)

func TestNewPermissionTableJson(t *testing.T) {

	testCases := []struct {
		jsonString    string
		expectedError string
		isNil         bool
		expectedRoles map[int]roles.Role
	}{
		{
			jsonString: `{"roles": {
				"1": {
					"name": "role 1",
					"permissions": [
						"perm1"
					]
				}
			}}`,
			expectedError: "",
			isNil:         false,
			expectedRoles: map[int]roles.Role{
				1: roles.Role{
					Name:        "role 1",
					Permissions: []string{"perm1"},
				},
			},
		},
		{
			jsonString: `{"roles": {
				"1": {
					"name": "role 1",
					"permissions": [
						"perm1"
					]
				},
				"12": {
					"name": "role 2",
					"permissions": [
						"perm1","perm3"
					]
				}
			}}`,
			expectedError: "",
			isNil:         false,
			expectedRoles: map[int]roles.Role{
				1: roles.Role{
					Name:        "role 1",
					Permissions: []string{"perm1"},
				},
				12: {
					"role 2",
					[]string{"perm1", "perm3"},
				},
			},
		},
		{
			jsonString: `{"roles": {
				"1": {
					"name": "role 1",
					"permissions": [
						"perm1"
					]
				},
				"alpha": {
					"name": "role 2",
					"permissions": [
						"perm1","perm3"
					]
				}
			}}`,
			expectedError: "json: cannot unmarshal number alpha into Go struct field PermissionTableJson.Roles of type int",
			isNil:         true,
			expectedRoles: map[int]roles.Role{},
		},
		{
			jsonString:    `{{}`,
			expectedError: "invalid character '{' looking for beginning of object key string",
			isNil:         true,
			expectedRoles: map[int]roles.Role{},
		},
	}

	for _, testCase := range testCases {
		// arrange
		file, err := ioutil.TempFile(os.TempDir(), "testJson*.json")
		if err != nil {
			assert.Nil(t, err)
		}
		file.WriteString(testCase.jsonString)
		file.Close()

		// act
		permTable, err := roles.NewPermissionTableJson(file.Name())

		// assert
		if len(testCase.expectedError) == 0 {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, testCase.expectedError, err.Error())
		}
		if testCase.isNil {
			assert.Nil(t, permTable)
		} else {
			assert.NotNil(t, permTable)
			assert.NotNil(t, permTable.Roles)
			assert.Equal(t, testCase.expectedRoles, permTable.Roles)
		}

		// cleanup
		os.Remove(file.Name())
	}
}

func TestPermissionTableJson_GetPermissions(t *testing.T) {
	testCases := []struct {
		table               *roles.PermissionTableJson
		roleId              int
		expectedPermissions []string
	}{
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
				},
			},
			roleId:              1,
			expectedPermissions: []string{"perm1"},
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleId:              12,
			expectedPermissions: []string{"perm3", "perm1"},
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleId:              1,
			expectedPermissions: []string{"perm1"},
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleId:              0,
			expectedPermissions: []string{},
		},
	}

	for _, testCase := range testCases {
		// arrange
		// act
		permissions := testCase.table.GetPermissions(testCase.roleId)

		// assert
		assert.Equal(t, testCase.expectedPermissions, permissions)
	}
}

func TestPermissionTableJson_LookUpRoleId(t *testing.T) {
	testCases := []struct {
		table          *roles.PermissionTableJson
		roleLookUp     string
		expectedRoleId int
		expectedError  error
	}{
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
				},
			},
			roleLookUp:     "role 1",
			expectedRoleId: 1,
			expectedError:  nil,
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleLookUp:     "role 12",
			expectedRoleId: 12,
			expectedError:  nil,
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleLookUp:     "role 1",
			expectedRoleId: 1,
			expectedError:  nil,
		},
		{
			table: &roles.PermissionTableJson{
				Roles: map[int]roles.Role{
					1: roles.Role{
						Name:        "role 1",
						Permissions: []string{"perm1"},
					},
					12: roles.Role{
						Name:        "role 12",
						Permissions: []string{"perm3", "perm1"},
					},
				},
			},
			roleLookUp:     "role 234",
			expectedRoleId: -1,
			expectedError:  errors.New("could not find role role 234"),
		},
	}

	for _, testCase := range testCases {
		// arrange
		// act
		id, err := testCase.table.LookUpRoleId(testCase.roleLookUp)

		// assert
		assert.Equal(t, testCase.expectedRoleId, id)
		assert.Equal(t, testCase.expectedError, err)
	}
}
