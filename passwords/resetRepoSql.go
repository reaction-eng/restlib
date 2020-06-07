// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package passwords

import (
	"database/sql"
	"errors"
	"time"

	"github.com/reaction-eng/restlib/configuration"
	"github.com/reaction-eng/restlib/email"
)

const TableName = "resetrequests"

var TokenExpired = errors.New("token_expired")

type ResetRepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//We need the emailer
	emailer                    email.Emailer
	resetEmailConfig           PasswordResetConfig
	activationEmailConfig      PasswordResetConfig
	oneTimePasswordEmailConfig PasswordResetConfig

	//Store the required statements to reduce compute time
	addRequestStatement *sql.Stmt
	getRequestStatement *sql.Stmt
	rmRequestStatement  *sql.Stmt

	tokenLifeSpan float64
}

/**
Store the type of token
*/
type tokenType int

const (
	activation      tokenType = 1
	reset           tokenType = 2
	oneTimePassword tokenType = 3
)

func NewRepoMySql(db *sql.DB, emailer email.Emailer, configuration configuration.Configuration) (*ResetRepoSql, error) {

	//Build a reset and activation config
	resetEmailConfig := PasswordResetConfig{}
	activationEmailConfig := PasswordResetConfig{}
	oneTimePasswordEmailConfig := PasswordResetConfig{}

	//Pull from the config
	err := configuration.GetStruct("password_reset", &resetEmailConfig)
	if err != nil {
		return nil, err
	}
	configuration.GetStruct("user_activation", &activationEmailConfig)
	if err != nil {
		return nil, err
	}
	configuration.GetStruct("one_time_password", &oneTimePasswordEmailConfig)
	if err != nil {
		return nil, err
	}
	tokenLifeSpan, err := configuration.GetFloat("tokenLifeSpan")
	if err != nil {
		return nil, err
	}

	//Define a new repo
	newRepo := ResetRepoSql{
		db:                         db,
		emailer:                    emailer,
		resetEmailConfig:           resetEmailConfig,
		activationEmailConfig:      activationEmailConfig,
		oneTimePasswordEmailConfig: oneTimePasswordEmailConfig,
		tokenLifeSpan:              tokenLifeSpan,
	}

	//Add request data to table
	addRequest, err := db.Prepare("INSERT INTO " + TableName + " (userId,email, token, issued, type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.addRequestStatement = addRequest

	//pull the request from the table
	getRequest, err := db.Prepare("SELECT * FROM " + TableName + " where userId = ? AND token = ? AND type = ?")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.getRequestStatement = getRequest

	//pull the request from the table
	rmRequest, err := db.Prepare("delete FROM " + TableName + " where id = ? limit 1")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.rmRequestStatement = rmRequest

	//Return a point
	return &newRepo, nil

}

func NewRepoPostgresSql(db *sql.DB, emailer email.Emailer, configuration configuration.Configuration) (*ResetRepoSql, error) {
	//Build a reset and activation config
	resetEmailConfig := PasswordResetConfig{}
	activationEmailConfig := PasswordResetConfig{}
	oneTimePasswordEmailConfig := PasswordResetConfig{}

	//Pull from the config
	err := configuration.GetStruct("password_reset", &resetEmailConfig)
	if err != nil {
		return nil, err
	}
	configuration.GetStruct("user_activation", &activationEmailConfig)
	if err != nil {
		return nil, err
	}
	configuration.GetStruct("one_time_password", &oneTimePasswordEmailConfig)
	if err != nil {
		return nil, err
	}
	tokenLifeSpan, err := configuration.GetFloat("tokenLifeSpan")
	if err != nil {
		return nil, err
	}

	//Define a new repo
	newRepo := ResetRepoSql{
		db:                         db,
		emailer:                    emailer,
		resetEmailConfig:           resetEmailConfig,
		activationEmailConfig:      activationEmailConfig,
		oneTimePasswordEmailConfig: oneTimePasswordEmailConfig,
		tokenLifeSpan:              tokenLifeSpan,
	}

	//Add request data to table
	addRequest, err := db.Prepare("INSERT INTO " + TableName + "(userId,email, token, issued, type) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.addRequestStatement = addRequest

	//pull the request from the table
	getRequest, err := db.Prepare("SELECT * FROM " + TableName + " where userId = $1 AND token = $2 AND type = $3")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.getRequestStatement = getRequest

	//pull the request from the table
	rmRequest, err := db.Prepare("delete FROM " + TableName + " where id = $1")
	if err != nil {
		return nil, err
	}

	//Store it
	newRepo.rmRequestStatement = rmRequest

	//Return a point
	return &newRepo, nil

}

func (repo *ResetRepoSql) IssueResetRequest(token string, userId int, emailAddress string) error {

	//Now add it to the database
	//execute the statement//(userId,name,input,flow)- "(userId,email, token, issued)
	_, err := repo.addRequestStatement.Exec(userId, emailAddress, token, time.Now(), reset)
	if err != nil {
		return err
	}

	//Make the email header
	header := email.HeaderInfo{
		Subject: repo.resetEmailConfig.Subject,
		To:      []string{emailAddress},
	}

	//Build a reset token
	resetInfo := PasswordResetInfo{
		Token: token,
		Email: emailAddress,
	}

	//Now email
	err = repo.emailer.SendTemplateFile(&header, repo.resetEmailConfig.Template, resetInfo, nil)

	return err
}

func (repo *ResetRepoSql) IssueActivationRequest(token string, userId int, emailAddress string) error {

	//Now add it to the database
	//execute the statement//(userId,name,input,flow)- "(userId,email, token, issued)
	_, err := repo.addRequestStatement.Exec(userId, emailAddress, token, time.Now(), activation)
	if err != nil {
		return err
	}

	//Make the email header
	header := email.HeaderInfo{
		Subject: repo.activationEmailConfig.Subject,
		To:      []string{emailAddress},
	}

	//Build a reset token
	resetInfo := PasswordResetInfo{
		Token: token,
		Email: emailAddress,
	}

	//Now email
	err = repo.emailer.SendTemplateFile(&header, repo.activationEmailConfig.Template, resetInfo, nil)

	return err
}

func (repo *ResetRepoSql) IssueOneTimePasswordRequest(token string, userId int, emailAddress string) error {

	//Now add it to the database
	_, err := repo.addRequestStatement.Exec(userId, emailAddress, token, time.Now(), oneTimePassword)
	if err != nil {
		return err
	}

	//Make the email header
	header := email.HeaderInfo{
		Subject: repo.activationEmailConfig.Subject,
		To:      []string{emailAddress},
	}

	//Build a reset token
	resetInfo := PasswordResetInfo{
		Token: token,
		Email: emailAddress,
	}

	//Now email
	err = repo.emailer.SendTemplateFile(&header, repo.oneTimePasswordEmailConfig.Template, resetInfo, nil)

	return err
}

func (repo *ResetRepoSql) CheckForResetToken(userId int, token string) (int, error) {

	//Get the id and errors
	id, err := repo.checkForToken(userId, token, reset)

	//If there is an error customize it
	if err != nil && err != TokenExpired {
		err = errors.New("password_change_forbidden")
	}

	return id, err

}

func (repo *ResetRepoSql) CheckForActivationToken(userId int, token string) (int, error) {

	//Get the id and errors
	id, err := repo.checkForToken(userId, token, activation)

	//If there is an error customize it
	if err != nil && err != TokenExpired {
		err = errors.New("activation_forbidden")
	}

	return id, err

}

func (repo *ResetRepoSql) CheckForOneTimePasswordToken(userId int, token string) (int, error) {

	//Get the id and errors
	id, err := repo.checkForToken(userId, token, oneTimePassword)

	//If there is an error customize it
	if err != nil && err != TokenExpired {
		err = errors.New("oneTimePassword_login_forbidden")
	}

	return id, err

}

func (repo *ResetRepoSql) checkForToken(userId int, token string, tkType tokenType) (int, error) {

	//Prepare to get values
	//id,  userId int, email TEXT, token TEXT, issued DATE,
	var id int
	var userIdDb int
	var emailDb string
	var tokenDb string
	var issued time.Time
	var tokenType tokenType

	//Get the value
	err := repo.getRequestStatement.QueryRow(userId, token, tkType).Scan(&id, &userIdDb, &emailDb, &tokenDb, &issued, &tokenType)

	//If there is an error, assume it can't be done
	if err != nil {
		return -1, errors.New("invalid_token")
	}

	//So it was correct, check the date
	if time.Now().Sub(issued).Hours() > repo.tokenLifeSpan {
		return 0, TokenExpired
	}

	//Make sure the user id and token match
	if userId != userIdDb || tokenDb != token {
		return -1, errors.New("invalid_token")
	}

	return id, nil
}

func (repo *ResetRepoSql) UseToken(id int) error {

	//Remove the token
	_, err := repo.rmRequestStatement.Exec(id)
	return err
}

func (repo *ResetRepoSql) CleanUp() {
	repo.getRequestStatement.Close()
	repo.addRequestStatement.Close()
	repo.rmRequestStatement.Close()
}
