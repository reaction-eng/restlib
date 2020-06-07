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

func (helper *BasicHelper) CreateUser(user User) error {

	//Make sure the info being passed in is valid
	if err := helper.validateUser(user); err != nil {
		return err
	}

	//Now hash the password
	user.SetPassword(helper.HashPassword(user.Password()))

	//Now store it
	newUser, err := helper.AddUser(user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Add the users to the org
	for _, orgId := range user.Organizations() {
		err = helper.AddUserToOrganization(newUser, orgId)
		if err != nil {
			return err
		}
	}

	//Else issue the request
	err = helper.IssueActivationRequest(helper.TokenGenerator(), newUser.Id(), newUser.Email())

	return err
}

/**
Validate incoming user details to make sure it has an email address and stuff,
*/
func (helper *BasicHelper) validateUser(user User) error {

	if !strings.Contains(user.Email(), "@") {
		return errors.New("validate_missing_email")
	}

	return helper.ValidatePassword(user.Password())
}

/**
Updates everything but the password
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

	//Don't add new orgs
	if !equal(newUser.Organizations(), oldUser.Organizations()) {
		return nil, errors.New("update_forbidden")
	}

	//Now update in the repo
	newUser, err = helper.UpdateUser(newUser)

	return newUser, err
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/**
Updates everything from the password
*/
func (helper *BasicHelper) PasswordChange(userId int, email string, newPassword string, oldPassword string) error {

	//Clean up the email
	email = strings.TrimSpace(strings.ToLower(email))

	//Load up the user
	oldUser, err := helper.GetUser(userId)
	if err != nil {
		return err
	}

	//Make sure the user can login with password
	if !oldUser.PasswordLogin() {
		return errors.New("user_password_login_forbidden")
	}

	//Make sure that the emails match
	if email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the old password matches
	passwordsMatch := helper.ComparePasswords(oldUser.Password(), oldPassword)

	//Make sure that the emails match
	if !passwordsMatch {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = helper.ValidatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.HashPassword(newPassword))

	//Now update in the repo
	_, err = helper.UpdateUser(oldUser)

	return err

}

func (helper *BasicHelper) PasswordChangeForced(userId int, email string, newPassword string) error {

	//Clean up the email
	email = strings.TrimSpace(strings.ToLower(email))

	//Load up the user
	oldUser, err := helper.GetUser(userId)
	if err != nil {
		return err
	}

	//Make sure that the emails match
	if email != oldUser.Email() {
		return errors.New("password_change_forbidden")
	}

	//Make sure the new password is valid
	err = helper.ValidatePassword(newPassword)

	//If the password is bad
	if err != nil {
		return err
	}

	//So it looks like we can update it, so hash the new password
	oldUser.SetPassword(helper.HashPassword(newPassword))

	//Now update in the repo
	_, err = helper.UpdateUser(oldUser)

	return err

}

/**
Login in the user
*/
func (helper *BasicHelper) Login(userPassword string, organizationId int, user User) (User, error) {

	//Make sure the user can login with password
	if !user.PasswordLogin() {
		return nil, errors.New("user_password_login_forbidden")
	}

	//Before you can login the user must be active
	if !user.Activated() {
		return nil, errors.New("user_not_activated")
	}

	// make sure user is in org
	if !InOrganization(user, organizationId) {
		return nil, errors.New("user_not_in_organization")
	}

	err := helper.ValidatePassword(userPassword)

	//If the password is bad
	if err != nil {
		return nil, errors.New("login_invalid_password")
	}

	//Now see if we login
	passwordsMatch := helper.ComparePasswords(user.Password(), userPassword)

	//Blank out the password before returning
	user.SetPassword("")

	//If they do not match
	if !passwordsMatch {
		return nil, errors.New("login_invalid_password")
	}

	//Create JWT token and Store the token in the response
	user.SetToken(helper.CreateJWTToken(user.Id(), organizationId, user.Email()))

	return user, nil
}
