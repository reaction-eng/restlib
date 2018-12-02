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
