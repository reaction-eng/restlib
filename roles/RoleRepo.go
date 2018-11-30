package roles

/**
Define an interface for roles
*/
type RoleRepo interface {
	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetPermissions(roleId int) []string
}
