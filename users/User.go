// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

//go:generate mockgen -destination=../mocks/mock_user.go -package=mocks github.com/reaction-eng/restlib/users User

type User interface {
	//Return the user id
	Id() int
	SetId(id int)

	//Return the user email
	Email() string
	SetEmail(email string)

	//Get the password
	Password() string
	SetPassword(password string)

	Token() string
	SetToken(token string)

	SetOrganizations(organizations ...int)
	Organizations() []int

	//Check if the user was activated
	Activated() bool

	//Check to see if the user can login with a password
	PasswordLogin() bool
}

func InOrganization(user User, organization int) bool {
	for _, org := range user.Organizations() {
		if org == organization {
			return true
		}
	}
	return false
}
