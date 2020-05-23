// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

//go:generate mockgen -destination=../mocks/mock_roles_repo.go -package=mocks -mock_names Repo=MockRolesRepo github.com/reaction-eng/restlib/roles  Repo

import "github.com/reaction-eng/restlib/users"

type Repo interface {
	GetPermissions(user users.User, organizationId int) (*Permissions, error)

	SetRolesByRoleId(user users.User, organizationId int, roles []int) error

	SetRolesByName(user users.User, organizationId int, roles []string) error
}
