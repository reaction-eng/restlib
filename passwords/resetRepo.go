// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords

//go:generate mockgen -destination=../mocks/mock_resetRepo.go -package=mocks github.com/reaction-eng/restlib/passwords ResetRepo

type ResetRepo interface {
	IssueResetRequest(token string, userId int, email string) error

	CheckForResetToken(userId int, resetToken string) (int, error)

	IssueActivationRequest(token string, userId int, email string) error

	CheckForActivationToken(userId int, activationToken string) (int, error)

	IssueOneTimePasswordRequest(token string, userId int, email string) error

	CheckForOneTimePasswordToken(userId int, activationToken string) (int, error)

	UseToken(id int) error

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
