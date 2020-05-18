// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords

//go:generate mockgen -destination=../mocks/mock_helper.go -package=mocks github.com/reaction-eng/restlib/passwords Helper

type Helper interface {
	HashPassword(password string) string
	CreateJWTToken(userId int, orgId int, email string) string
	ComparePasswords(currentPwHash string, testingPassword string) bool
	TokenGenerator() string
	ValidateToken(tokenHeader string) (int, int, string, error)
	ValidatePassword(password string) error
}
