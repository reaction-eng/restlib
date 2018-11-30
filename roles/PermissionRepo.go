package roles

import "bitbucket.org/reidev/restlib/users"

/**
Define an interface for roles
*/
type PermissionRepo interface {
	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetPermissions(user users.User) (*Permissions, error)
}
