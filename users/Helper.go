package users

import (
	"bitbucket.org/reidev/restlib/passwords"
	"errors"
	"strings"
)

type Helper struct {

	//Hold the user repo
	usersRepo Repo

	//And the password repo
	passRepo passwords.ResetRepo

	//And store a password helper
	passwordHelper passwords.Helper
}

func NewUserHelper(usersRepo Repo, passRepo passwords.ResetRepo, passwordHelper passwords.Helper) *Helper {

	return &Helper{
		usersRepo:      usersRepo,
		passRepo:       passRepo,
		passwordHelper: passwordHelper,
	}

}

/**
Static method to create a new user
*/
func (helper *Helper) createUser(user User) error {

	//Make sure the info being passed in is valid
	if ok, err := helper.validateUser(user); !ok {
		return err
	}

	//Now hash the password
	user.SetPassword(helper.passwordHelper.HashPassword(user.Password()))

	//Now store it
	newUser, err := helper.usersRepo.AddUser(user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Else issue the request
	err = helper.passRepo.IssueActivationRequest(helper.passwordHelper.TokenGenerator(), newUser.Id(), newUser.Email())

	if err != nil {
		return err
	}

	return nil

}

/**
Validate incoming user details to make sure it has an email address and stuff
*/
func (helper *Helper) validateUser(user User) (bool, error) {

	if !strings.Contains(user.Email(), "@") {
		return false, errors.New("validate_missing_email")
	}

	//Check the password
	err := helper.passwordHelper.ValidatePassword(user.Password())

	//If the user already exists
	if err != nil {
		return false, err
	}

	//Now look up a possible user
	user, err = helper.usersRepo.GetUserByEmail(user.Email())

	//If the user already exists
	if err == nil || user != nil {
		return false, errors.New("validate_email_in_use")
	}

	//All is good
	return true, nil
}

/**
Updates everything from the password
*/
func (helper *Helper) updateUser(userId int, newUser User) (User, error) {

	//Load up the user
	oldUser, err := helper.usersRepo.GetUser(userId)

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
	newUser, err = helper.usersRepo.UpdateUser(newUser)

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
func (helper *Helper) passwordChange(userId int, passwordChange updatePasswordChangeStruct) error {

	//Clean up the email
	passwordChange.Email = strings.TrimSpace(strings.ToLower(passwordChange.Email))

	//Load up the user
	oldUser, err := helper.usersRepo.GetUser(userId)

	//Make sure the user can login with password
	if !oldUser.PasswordLogin() {
		return errors.New("user_password_login_forbidden")
	}

	//Make sure that the emails match
	if passwordChange.Email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the old password matches
	passwordsMath := helper.passwordHelper.ComparePasswords(oldUser.Password(), passwordChange.PasswordOld)

	//Make sure that the emails match
	if !passwordsMath {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = helper.passwordHelper.ValidatePassword(passwordChange.Password)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.passwordHelper.HashPassword(passwordChange.Password))

	//Now update in the repo
	_, err = helper.usersRepo.UpdateUser(oldUser)

	return err

}

/**
Updates everything from the password
*/
func (helper *Helper) passwordChangeForced(userId int, email string, newPassword string) error {

	//Clean up the email
	email = strings.TrimSpace(strings.ToLower(email))

	//Load up the user
	oldUser, err := helper.usersRepo.GetUser(userId)

	//Make sure the user can login with password
	//if !oldUser.PasswordLogin() {
	//	return errors.New("user_password_login_forbidden")
	//}

	//Make sure the new password is valid
	err = helper.passwordHelper.ValidatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.passwordHelper.HashPassword(newPassword))

	//Now update in the repo
	_, err = helper.usersRepo.UpdateUser(oldUser)

	return err

}

/**
Login in the user
*/
func (helper *Helper) login(userPassword string, user User) (User, error) {

	//Make sure the user can login with password
	if !user.PasswordLogin() {
		return nil, errors.New("user_password_login_forbidden")
	}

	//Before you can login the user must be active
	if !user.Activated() {
		return nil, errors.New("user_not_activated")
	}

	//Make sure the new password is valid
	err := helper.passwordHelper.ValidatePassword(userPassword)

	//If the password is bad
	if err != nil {
		return nil, errors.New("login_invalid_password")
	}

	//Now see if we login
	passwordsMath := helper.passwordHelper.ComparePasswords(user.Password(), userPassword)

	//Blank out the password before returning
	user.SetPassword("")

	//If they do not match
	if !passwordsMath {
		return nil, errors.New("login_invalid_password")
	}

	//Create JWT token and Store the token in the response
	user.SetToken(helper.passwordHelper.CreateJWTToken(user.Id(), user.Email()))

	return user, nil
}
