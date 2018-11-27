package users

//a struct to rep user account
type User interface {
	//Return the user id
	Id() int
	SetId(id int)

	//Return the user email
	Email() string

	//Get the password
	Password() string
	SetPassword(password string)

	//Return the user email
	Token() string
	SetToken(token string)
}
