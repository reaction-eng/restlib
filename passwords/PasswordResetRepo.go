package passwords

/**
Define an interface that all Calc Repos must follow
*/
type PasswordResetRepo interface {

	/**
	Issues a request for the user
	*/
	IssueResetRequest(userId int, email string) error

	/**
	Issues a request for the user
	*/
	CheckForResetToken(userId int, reset_token string) (int, error)

	/**
	Issues a request for the user
	*/
	UseResetToken(id int) error

	/**
	Allow databases to be closed
	*/
	CleanUp()
}
