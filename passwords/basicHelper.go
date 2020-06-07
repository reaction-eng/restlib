// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/reaction-eng/restlib/configuration"
	"golang.org/x/crypto/bcrypt"
)

type BasicHelper struct {
	//Keep a global password config
	jwtTokenPassword []byte
}

/*
JWT claims struct
*/
type Token struct {
	UserId         int
	OrganizationId int
	Email          string
	jwt.StandardClaims
}

//Load it during init
func NewBasicHelper(configuration configuration.Configuration) (*BasicHelper, error) {
	//Now get the token
	jwtTokenPasswordString := configuration.GetString("token_password")

	//If it is null error
	if len(jwtTokenPasswordString) < 60 {
		return nil, errors.New("the jwt token is not specified or not long enough")
	}

	//Store the byte array
	return &BasicHelper{
		jwtTokenPassword: []byte(jwtTokenPasswordString),
	}, nil
}

/**
Support function to hash the password
*/
func (helper *BasicHelper) HashPassword(password string) string {

	//Hash the password, there should be a salt
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

/**
  Support function to generate a JWT token
*/
func (helper *BasicHelper) CreateJWTToken(userId int, organizationId int, email string) string {

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: userId, OrganizationId: organizationId, Email: email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString(helper.jwtTokenPassword)

	return tokenString

}

/**
  Compare mysql.  Determine if they match
*/
func (helper *BasicHelper) ComparePasswords(currentPwHash string, testingPassword string) bool {

	//Now take the password and encrypt it
	err := bcrypt.CompareHashAndPassword([]byte(currentPwHash), []byte(testingPassword))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false
	} else {
		return true
	}
}

/**
 * Get a random token
 */
func (helper *BasicHelper) TokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (helper *BasicHelper) ValidateToken(tokenHeader string) (int, int, string, error) {

	//Token is missing, returns with error code 403 Unauthorized
	if tokenHeader == "" {
		return -1, -1, "", errors.New("auth_missing_token")
	}

	//Now split the token to get the useful part
	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		return -1, -1, "", errors.New("auth_malformed_token")

	}

	//Grab the token part, what we are truly interested in
	tokenPart := splitted[1]

	//Get the token and take it back apart
	tk := &Token{}

	//Now parse the token
	token, err := jwt.ParseWithClaims(tokenPart, tk,
		func(token *jwt.Token) (interface{}, error) {
			return helper.jwtTokenPassword, nil
		})

	//check for mailformed data
	if err != nil { //Malformed token, returns with http code 403 as usual
		return -1, -1, "", errors.New("auth_malformed_token")

	}

	//Token is invalid, maybe not signed on this server
	if !token.Valid {
		//Return the error
		return -1, -1, "", errors.New("auth_forbidden")
	}

	return tk.UserId, tk.OrganizationId, tk.Email, nil
}

/**
Make sure that the password is valid
*/
func (helper *BasicHelper) ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("validate_password_insufficient")
	}
	return nil
}
