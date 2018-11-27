package users

import (
	"bitbucket.org/reidev/restlib/authentication"
	"errors"
	"strings"
)

/**
Static method to create a new user
*/
func CreateUser(usersRepo Repo, user User) error {

	//Make sure the info being passed in is valid
	if ok, err := ValidateUser(usersRepo, user); !ok {
		return err
	}

	//Now hash the password
	user.SetPassword(authentication.HashPassword(user.Password()))

	//Now store it
	userReturn, err := usersRepo.AddUser(user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Store the return
	user = userReturn

	//Store the token to return
	user.SetToken(authentication.CreateJWTToken(user.Id()))

	//Clear the password
	user.SetPassword("") //delete password

	return nil

}

/**
Validate incoming user details to make sure it has an email address and stuff
*/
func ValidateUser(usersRepo Repo, user User) (bool, error) {

	if !strings.Contains(user.Email(), "@") {
		return false, errors.New("validate_missing_email")
	}

	if len(user.Password()) < 6 {
		return false, errors.New("validate_password_insufficient")
	}

	//Now look up a possible user
	_, err := usersRepo.GetUserByEmail(user.Email())

	//If the user already exists
	if err == nil {
		return false, errors.New("validate_email_in_use")
	}

	//All is good
	return true, nil
}

/**
Login in the user
*/
func Login(userPassword string, user User) error {

	//Now see if we login
	passwordsMath := authentication.ComparePasswords(user.Password(), userPassword)

	//Blank out the password before returning
	user.SetPassword("")

	//If they do not match
	if !passwordsMath {
		return errors.New("login_invalid_password")
	}

	//Create JWT token and Store the token in the response
	user.SetToken(authentication.CreateJWTToken(user.Id()))

	return nil
}
