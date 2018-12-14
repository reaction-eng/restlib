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
	IssueActivationRequest(userId int, email string) error

	/**
	Issues a request for the user
	*/
	CheckForActivationToken(userId int, activationToken string) (int, error)

	/**
	Issues a request for the user
	*/
	UseToken(id int) error

	/**
	Allow databases to be closed
	*/
	CleanUp()
}

//Define a struct to store password reset configs
type PasswordResetConfig struct {
	Template string `json:"template"`
	Subject  string `json:"subject"`
}

//Define a struct to store password reset configs
type PasswordResetInfo struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
