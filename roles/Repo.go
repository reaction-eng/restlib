package roles

import "bitbucket.org/reidev/restlib/users"

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
