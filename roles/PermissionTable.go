// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

/**
Define an interface for roles
*/
type PermissionTable interface {
	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetPermissions(roleId int) []string

	//Look up the role id based upon the name
	LookUpRoleId(name string) (int, error)
}
