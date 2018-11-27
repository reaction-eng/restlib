package authentication

/**
Define an interface that all Calc Repos must follow
*/
type PasswordResetRepo interface {

	/**
	Issues a request for the user
	*/
	IssueResetRequest(userId int, email string) error

	/**
	Allow databases to be closed
	*/
	CleanUp()
}
