package users

import (
	"bitbucket.org/reidev/restlib/passwords"
	"errors"
	"strings"
)

/**
Static method to create a new user
*/
func createUser(usersRepo Repo, passRepo passwords.PasswordResetRepo, user User) error {

	//Make sure the info being passed in is valid
	if ok, err := validateUser(usersRepo, user); !ok {
		return err
	}

	//Now hash the password
	user.SetPassword(passwords.HashPassword(user.Password()))

	//Now store it
	newUser, err := usersRepo.AddUser(user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Else issue the request
	err = passRepo.IssueActivationRequest(newUser.Id(), newUser.Email())

	if err != nil {
		return err
	}

	return nil

}

/**
Validate incoming user details to make sure it has an email address and stuff
*/
func validateUser(usersRepo Repo, user User) (bool, error) {

	if !strings.Contains(user.Email(), "@") {
		return false, errors.New("validate_missing_email")
	}

	//Check the password
	err := validatePassword(user.Password())

	//If the user already exists
	if err != nil {
		return false, err
	}

	//Now look up a possible user
	user, err = usersRepo.GetUserByEmail(user.Email())

	//If the user already exists
	if err == nil || user != nil {
		return false, errors.New("validate_email_in_use")
	}

	//All is good
	return true, nil
}

/**
Make sure that the password is valid
*/
func validatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("validate_password_insufficient")
	}
	return nil
}

/**
Updates everything from the password
*/
func updateUser(usersRepo Repo, userId int, newUser User) (User, error) {

	//Load up the user
	oldUser, err := usersRepo.GetUser(userId)

	//Check for err
	if err != nil {
		return nil, err
	}

	//There are three things we cannot change when we update the user, the id
	if newUser.Id() != oldUser.Id() {
		return nil, errors.New("update_forbidden")
	}

	//And the password
	if newUser.Password() != oldUser.Password() {
		return nil, errors.New("update_forbidden")
	}

	//And the email
	if newUser.Email() != oldUser.Email() {
		return nil, errors.New("update_forbidden")
	}

	//Make sure we

	//Now update in the repo
	newUser, err = usersRepo.UpdateUser(newUser)

	return newUser, err

}

/**
Define a struct for just updating password
*/
type updatePasswordChangeStruct struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PasswordOld string `json:"passwordold"`
}

/**
Updates everything from the password
*/
func passwordChange(usersRepo Repo, userId int, passwordChange updatePasswordChangeStruct) error {

	//Load up the user
	oldUser, err := usersRepo.GetUser(userId)

	//Make sure that the emails match
	if passwordChange.Email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the old password matches
	passwordsMath := passwords.ComparePasswords(oldUser.Password(), passwordChange.PasswordOld)

	//Make sure that the emails match
	if !passwordsMath {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = validatePassword(passwordChange.Password)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(passwords.HashPassword(passwordChange.Password))

	//Now update in the repo
	_, err = usersRepo.UpdateUser(oldUser)

	return err

}

/**
Updates everything from the password
*/
func passwordChangeForced(usersRepo Repo, userId int, email string, newPassword string) error {

	//Load up the user
	oldUser, err := usersRepo.GetUser(userId)

	//Make sure that the emails match
	if email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = validatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(passwords.HashPassword(newPassword))

	//Now update in the repo
	_, err = usersRepo.UpdateUser(oldUser)

	return err

}

/**
Login in the user
*/
func login(userPassword string, user User) (User, error) {

	//Before you can login the user must be active
	if !user.Activated() {
		return nil, errors.New("user_not_activated")
	}

	//Make sure the new password is valid
	err := validatePassword(userPassword)

	//If the password is bad
	if err != nil {
		return nil, errors.New("login_invalid_password")
	}

	//Now see if we login
	passwordsMath := passwords.ComparePasswords(user.Password(), userPassword)

	//Blank out the password before returning
	user.SetPassword("")

	//If they do not match
	if !passwordsMath {
		return nil, errors.New("login_invalid_password")
	}

	//Create JWT token and Store the token in the response
	user.SetToken(passwords.CreateJWTToken(user.Id(), user.Email()))

	return user, nil
}
