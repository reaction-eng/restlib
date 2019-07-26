// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords

/**
Define an interface that all Calc Repos must follow
*/
type ResetRepo interface {

	/**
	Issues a request for the user
	*/
	IssueResetRequest(token string, userId int, email string) error

	/**
	Issues a request for the user
	*/
	CheckForResetToken(userId int, reset_token string) (int, error)

	/**
	Issues a request for the user
	*/
	IssueActivationRequest(token string, userId int, email string) error

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
