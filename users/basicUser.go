// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import "strings"

type BasicUser struct {
	Id_            int    `json:"id"`
	Organizations_ []int  `json:"organizations"`
	Email_         string `json:"email"`
	password_      string `json:"-"`
	Token_         string `json:"token";sql:"-"`
	activated_     bool
	passwordlogin_ bool
}

/**
Add the required setters and getters
*/
func (basic *BasicUser) Id() int {
	return basic.Id_
}
func (basic *BasicUser) SetId(id int) {
	basic.Id_ = id
}
func (basic *BasicUser) Email() string {
	return strings.TrimSpace(strings.ToLower(basic.Email_))
}
func (basic *BasicUser) SetEmail(email string) {
	basic.Email_ = email
}

//func (basic *BasicUser) SetId(id int)  {
//	basic.Id_ = id
//}
func (basic *BasicUser) Password() string {
	return basic.password_
}
func (basic *BasicUser) SetPassword(pw string) {
	basic.password_ = pw
}
func (basic *BasicUser) Token() string {
	return basic.Token_
}
func (basic *BasicUser) SetToken(tk string) {
	basic.Token_ = tk
}

func (basic *BasicUser) Activated() bool {
	return basic.activated_
}

func (basic *BasicUser) PasswordLogin() bool {
	return basic.passwordlogin_
}

func (basic *BasicUser) Organizations() []int {
	return basic.Organizations_
}

func (basic *BasicUser) SetOrganizations(organizations ...int) {
	basic.Organizations_ = organizations
}

/**
Provide code to copy the user into this user
*/
/**
Add the required setters and getters
*/
func (basic *BasicUser) CopyFrom(from User) {
	basic.Email_ = from.Email()
	basic.password_ = from.Password()
	basic.Id_ = from.Id()
	basic.Token_ = from.Token()
	basic.activated_ = from.Activated()
	basic.passwordlogin_ = from.PasswordLogin()
	basic.Organizations_ = from.Organizations()
}
