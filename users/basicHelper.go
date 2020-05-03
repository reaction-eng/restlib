// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"errors"
	"strings"

	"github.com/reaction-eng/restlib/passwords"
)

type BasicHelper struct {
	//Hold the user repo
	Repo

	//And the password repo
	passwords.ResetRepo

	//And store a password helper
	passwords.Helper
}

func NewUserHelper(usersRepo Repo, passRepo passwords.ResetRepo, passwordHelper passwords.Helper) *BasicHelper {
	return &BasicHelper{
		Repo:      usersRepo,
		ResetRepo: passRepo,
		Helper:    passwordHelper,
	}
}

/**
Static method to create a new user
*/
func (helper *BasicHelper) CreateUser(user User) error {

	//Make sure the info being passed in is valid
	if ok, err := helper.ValidateUser(user); !ok {
		return err
	}

	//Now hash the password
	user.SetPassword(helper.Helper.HashPassword(user.Password()))

	//Now store it
	newUser, err := helper.AddUser(user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Else issue the request
	err = helper.IssueActivationRequest(helper.Helper.TokenGenerator(), newUser.Id(), newUser.Email())

	if err != nil {
		return err
	}

	return nil

}

/**
Validate incoming user details to make sure it has an email address and stuff,
//TODO: add organization check
*/
func (helper *BasicHelper) ValidateUser(user User) (bool, error) {

	if !strings.Contains(user.Email(), "@") {
		return false, errors.New("validate_missing_email")
	}

	//Check the password
	err := helper.Helper.ValidatePassword(user.Password())

	//If the user already exists
	if err != nil {
		return false, err
	}

	//Now look up a possible user
	user, err = helper.GetUserByEmail(user.Email())

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
func (helper *BasicHelper) Update(userId int, newUser User) (User, error) {

	//Load up the user
	oldUser, err := helper.GetUser(userId)

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
	newUser, err = helper.UpdateUser(newUser)

	return newUser, err

}

/**
Updates everything from the password
*/
func (helper *BasicHelper) PasswordChange(userId int, email string, newPassword string, oldPassword string) error {

	//Clean up the email
	email = strings.TrimSpace(strings.ToLower(email))

	//Load up the user
	oldUser, err := helper.GetUser(userId)

	//Make sure the user can login with password
	if !oldUser.PasswordLogin() {
		return errors.New("user_password_login_forbidden")
	}

	//Make sure that the emails match
	if email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the old password matches
	passwordsMath := helper.Helper.ComparePasswords(oldUser.Password(), oldPassword)

	//Make sure that the emails match
	if !passwordsMath {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = helper.Helper.ValidatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.Helper.HashPassword(newPassword))

	//Now update in the repo
	_, err = helper.UpdateUser(oldUser)

	return err

}

/**
Updates everything from the password
*/
func (helper *BasicHelper) PasswordChangeForced(userId int, email string, newPassword string) error {

	//Clean up the email
	email = strings.TrimSpace(strings.ToLower(email))

	//Load up the user
	oldUser, err := helper.GetUser(userId)

	//Make sure the user can login with password
	//if !oldUser.PasswordLogin() {
	//	return errors.New("user_password_login_forbidden")
	//}

	//Make sure the new password is valid
	err = helper.Helper.ValidatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.Helper.HashPassword(newPassword))

	//Now update in the repo
	_, err = helper.UpdateUser(oldUser)

	return err

}

/**
Login in the user
*/
func (helper *BasicHelper) Login(userPassword string, orgId int, user User) (User, error) {

	//Make sure the user can login with password
	if !user.PasswordLogin() {
		return nil, errors.New("user_password_login_forbidden")
	}

	//Before you can login the user must be active
	if !user.Activated() {
		return nil, errors.New("user_not_activated")
	}

	//Make sure the new password is valid
	err := helper.Helper.ValidatePassword(userPassword)

	//If the password is bad
	if err != nil {
		return nil, errors.New("login_invalid_password")
	}

	//Now see if we login
	passwordsMath := helper.Helper.ComparePasswords(user.Password(), userPassword)

	//Blank out the password before returning
	user.SetPassword("")

	//If they do not match
	if !passwordsMath {
		return nil, errors.New("login_invalid_password")
	}

	//Create JWT token and Store the token in the response
	user.SetToken(helper.Helper.CreateJWTToken(user.Id(), user.Email()))

	return user, nil
}
