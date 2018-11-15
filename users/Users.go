package users

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId int
	jwt.StandardClaims
}

//a struct to rep user account
type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (user *User) Create(usersRepo Repo) error {

	//Make sure the info being passed in is valid
	if ok, err := user.Validate(usersRepo); !ok {
		return err
	}

	//Hash the password, there should be a salt
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	//Now store it
	userReturn, err := usersRepo.AddUser(*user)

	//Make sure it created an id
	if err != nil {
		return err
	}

	//Store the return
	*user = userReturn

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	//Store the token to return
	user.Token = tokenString

	//Clear the password
	user.Password = "" //delete password

	return nil

}

/**
Validate incoming user details to make sure it has an email address and stuff
*/
func (user *User) Validate(usersRepo Repo) (bool, error) {

	if !strings.Contains(user.Email, "@") {
		return false, errors.New("email address is required")
	}

	if len(user.Password) < 6 {
		return false, errors.New("password does not meeting length requirements")
	}

	//Now look up a possible user
	_, err := usersRepo.GetUserByEmail(user.Email)

	//If the user already exisits
	if err == nil {
		return false, errors.New("email already in use")
	}

	//All is good
	return true, nil
}

/**
Login in the user
*/
func (user *User) Login(userPassword string) error {

	//Now take the password and encrypt it
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		user.Password = ""
		return errors.New("invalid login credentials. please try again")
	}

	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &Token{UserId: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	//Store the token in the response
	user.Token = tokenString

	return nil
}

//
//func GetUserByEmail(u uint) *Account {
//
//	acc := &Account{}
//	GetDB().Table("accounts").Where("id = ?", u).First(acc)
//	if acc.Email == "" { //User not found!
//		return nil
//	}
//
//	acc.Password = ""
//	return acc
//}
