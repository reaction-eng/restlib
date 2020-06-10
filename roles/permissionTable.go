// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

//go:generate mockgen -destination=../mocks/mock_permissionTable.go -package=mocks github.com/reaction-eng/restlib/roles PermissionTable

type PermissionTable interface {
	GetPermissions(roleId int) []string

	LookUpRoleId(name string) (int, error)
}
