package users

/**
Define an interface that all Calc Repos must follow
*/
type Repo interface {
	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetUserByEmail(email string) (User, error)

	/**
	Get the user with the email.  An error is thrown is not found
	*/
	GetUser(id int) (User, error)

	/**
	Add User
	*/
	AddUser(user User) (User, error)

	/**
	Update User
	*/
	UpdateUser(user User) (User, error)

	/**
	Allow databases to be closed
	*/
	CleanUp()

	/**
	Create empty user
	*/
	NewEmptyUser() User
}
