// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

//go:generate mockgen -destination=../mocks/mock_roles_repo.go -package=mocks -mock_names Repo=MockRolesRepo github.com/reaction-eng/restlib/roles  Repo

import "github.com/reaction-eng/restlib/users"

/**
Define an interface for roles
*/
type Repo interface {
	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetPermissions(user users.User) (*Permissions, error)

	/**
	Set the user's roles. Note this wipes out all current roles
	*/
	SetRolesByRoleId(user users.User, roles []int) error

	/**
	Set the user's roles.  Note this wipes out all current roles
	*/
	SetRolesByName(user users.User, roles []string) error
}
