// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/reaction-eng/restlib/utils"
)

const UserTableName = "users"
const UserOrgTableName = "userOrganizations"

/**
Define a struct for Repo for use with users
*/
type RepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//Also store the table name
	tableName string

	//Store the required statements to reduce compute time
	addUserStatement           *sql.Stmt
	getUserStatement           *sql.Stmt
	getUserByEmailStatement    *sql.Stmt
	updateUserStatement        *sql.Stmt
	activateStatement          *sql.Stmt
	listAllUsersStatement      *sql.Stmt
	getUserOrganizations       *sql.Stmt
	addUserToOrganization      *sql.Stmt
	removeUserFromOrganization *sql.Stmt
}

//Provide a method to make a new UserRepoSql
func NewRepoMySql(db *sql.DB) (*RepoSql, error) {

	//Define a new repo
	newRepo := RepoSql{
		db: db,
	}

	addUser, err := db.Prepare("INSERT INTO " + UserTableName + "(email,password) VALUES (?, ?)")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.addUserStatement = addUser

	getUser, err := db.Prepare("SELECT * FROM " + UserTableName + " where id = ?")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.getUserStatement = getUser

	getUserByEmail, err := db.Prepare("SELECT * FROM " + UserTableName + " where email like ?")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.getUserByEmailStatement = getUserByEmail

	updateStatement, err := db.Prepare("UPDATE  " + UserTableName + " SET email = ?, password = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	newRepo.updateUserStatement = updateStatement

	activateStatement, err := db.Prepare("UPDATE  " + UserTableName + " SET activation = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	newRepo.activateStatement = activateStatement

	listAllUsers, err := db.Prepare("SELECT id, activation FROM " + UserTableName)
	if err != nil {
		return nil, err
	}
	newRepo.listAllUsersStatement = listAllUsers

	getUserOrganizations, err := db.Prepare("SELECT orgId FROM " + UserOrgTableName + " where userId = ?")
	if err != nil {
		return nil, err
	}
	newRepo.getUserOrganizations = getUserOrganizations

	addUserToOrganization, err := db.Prepare("INSERT INTO " + UserOrgTableName + " (userId,orgId,joinDate) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	newRepo.addUserToOrganization = addUserToOrganization

	removeUserFromOrganization, err := db.Prepare("DELETE FROM " + UserOrgTableName + " WHERE userId = ? AND orgId = ?")
	if err != nil {
		return nil, err
	}
	newRepo.removeUserFromOrganization = removeUserFromOrganization

	//Return a point
	return &newRepo, nil

}

func NewRepoPostgresSql(db *sql.DB) (*RepoSql, error) {

	//Define a new repo
	newRepo := RepoSql{
		db: db,
	}

	addUser, err := db.Prepare("INSERT INTO " + UserTableName + "(email,password) VALUES ($1, $2)")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.addUserStatement = addUser

	getUser, err := db.Prepare("SELECT * FROM " + UserTableName + " where id = $1")
	if err != nil {
		return nil, err
	}
	newRepo.getUserStatement = getUser

	getUserByEmail, err := db.Prepare("SELECT * FROM " + UserTableName + " where email like $1")
	if err != nil {
		return nil, err
	}
	newRepo.getUserByEmailStatement = getUserByEmail

	updateStatement, err := db.Prepare("UPDATE  " + UserTableName + " SET email = $1, password = $2 WHERE id = $3")
	if err != nil {
		return nil, err
	}
	newRepo.updateUserStatement = updateStatement

	activateStatement, err := db.Prepare("UPDATE  " + UserTableName + " SET activation = $1 WHERE id = $2")
	if err != nil {
		return nil, err
	}
	newRepo.activateStatement = activateStatement

	listAllUsers, err := db.Prepare("SELECT id, activation FROM " + UserTableName)
	if err != nil {
		return nil, err
	}
	newRepo.listAllUsersStatement = listAllUsers

	getUserOrganizations, err := db.Prepare("SELECT orgId FROM " + UserOrgTableName + " where userId = $1")
	if err != nil {
		return nil, err
	}
	newRepo.getUserOrganizations = getUserOrganizations

	addUserToOrganization, err := db.Prepare("INSERT INTO " + UserOrgTableName + " (userId,orgId,joinDate) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	newRepo.addUserToOrganization = addUserToOrganization

	removeUserFromOrganization, err := db.Prepare("DELETE FROM " + UserOrgTableName + " WHERE userId = $1 AND orgId = $2")
	if err != nil {
		return nil, err
	}
	newRepo.removeUserFromOrganization = removeUserFromOrganization

	return &newRepo, nil
}

/**
Look up the user and return if they were found
*/
func (repo *RepoSql) GetUserByEmail(email string) (User, error) {
	//Clean up the string
	email = strings.TrimSpace(strings.ToLower(email))

	//var dataResult string
	var user BasicUser

	//Store the sql time
	var activationDate utils.NullTime

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserByEmailStatement.QueryRow(email).Scan(&user.Id_, &user.Email_, &user.password_, &activationDate)

	//Use a useful error
	if err == sql.ErrNoRows {
		return nil, UserNotFound
	} else if err != nil {
		return nil, err
	}

	//Store if this is activated
	user.activated_ = activationDate.Valid
	user.passwordlogin_ = len(user.password_) > 0

	orgIds, err := repo.listUserOrganizations(user.Id_)
	if err != nil {
		return nil, err
	}
	user.SetOrganizations(orgIds...)

	return &user, err
}

/**
Look up the user by id and return if they were found
*/
func (repo *RepoSql) GetUser(id int) (User, error) {
	var user BasicUser

	//Store the sql time
	var activationDate utils.NullTime

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserStatement.QueryRow(id).Scan(&user.Id_, &user.Email_, &user.password_, &activationDate)

	//Use a useful error
	if err == sql.ErrNoRows {
		return nil, UserNotFound
	} else if err != nil {
		return nil, err
	}

	//Store if this is activated
	user.activated_ = activationDate.Valid
	user.passwordlogin_ = len(user.password_) > 0

	orgIds, err := repo.listUserOrganizations(user.Id_)
	if err != nil {
		return nil, err
	}
	user.SetOrganizations(orgIds...)

	return &user, err
}

/**
List all of the users
*/
func (repo *RepoSql) ListUsers(onlyActive bool, organizations []int) ([]int, error) {
	//Put in the list
	list := make([]int, 0)

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	rows, err := repo.listAllUsersStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		//Store the sql time
		var activationDate utils.NullTime

		err := rows.Scan(&id, &activationDate)
		if err != nil {
			return nil, err
		}

		//Append the row
		if !onlyActive || activationDate.Valid {
			list = append(list, id)
		}
	}
	err = rows.Err()

	return list, err
}

func (repo *RepoSql) AddUser(newUser User) (User, error) {

	_, userFoundError := repo.GetUserByEmail(newUser.Email())
	if userFoundError == nil {
		return nil, errors.New("user_email_in_user")
	}
	if userFoundError != UserNotFound {
		return nil, userFoundError
	}

	//Add the info
	//execute the statement//(userId,name,input,flow)
	_, err := repo.addUserStatement.Exec(newUser.Email(), newUser.Password())

	//Check for error
	if err != nil {
		return nil, err
	}

	//Now look up the person by email
	return repo.GetUserByEmail(newUser.Email())

}

func (repo *RepoSql) UpdateUser(user User) (User, error) {
	//execute the statement//"UPDATE  " + tableName + " SET email = ?, password = ? WHERE id = ?"
	_, err := repo.updateUserStatement.Exec(user.Email(), user.Password(), user.Id())

	return user, err
}

func (repo *RepoSql) ActivateUser(user User) error {
	//Get the current time
	actTime := utils.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	//Just update the info//"UPDATE  " + tableName + " SET activation = $1 WHERE id = $2")
	_, err := repo.activateStatement.Exec(actTime, user.Id())
	return err
}

func (repo *RepoSql) AddUserToOrganization(user User, orgId int) error {
	_, err := repo.addUserToOrganization.Exec(user.Id(), orgId, time.Now())
	return err
}

func (repo *RepoSql) RemoveUserFromOrganization(user User, orgId int) error {
	result, err := repo.removeUserFromOrganization.Exec(user.Id(), orgId)
	if err != nil {
		return err
	}

	rowsChange, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsChange < 1 {
		return errors.New("no_organizations_removed")
	}

	return err
}

func (repo *RepoSql) listUserOrganizations(userId int) ([]int, error) {
	//Put in the list
	list := make([]int, 0)

	rows, err := repo.getUserOrganizations.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orgId int

		err := rows.Scan(&orgId)
		if err != nil {
			return nil, err
		}

		list = append(list, orgId)
	}
	err = rows.Err()

	return list, err

}

func (repo *RepoSql) CleanUp() {
	repo.addUserStatement.Close()
	repo.getUserByEmailStatement.Close()
	repo.getUserStatement.Close()
	repo.updateUserStatement.Close()
	repo.listAllUsersStatement.Close()
}

func (repo *RepoSql) NewEmptyUser() User {
	return &BasicUser{}
}
