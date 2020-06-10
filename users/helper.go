package users

//go:generate mockgen -destination=../mocks/mock_user_helper.go -package=mocks -mock_names Helper=MockUserHelper github.com/reaction-eng/restlib/users  Helper

import (
	"github.com/reaction-eng/restlib/passwords"
)

type Helper interface {
	Repo
	passwords.ResetRepo
	passwords.Helper

	CreateUser(user User) error
	Update(userId int, newUser User) (User, error)
	PasswordChange(userId int, email string, newPassword string, oldPassword string) error
	PasswordChangeForced(userId int, email string, newPassword string) error
	Login(userPassword string, orgId int, user User) (User, error)
}
