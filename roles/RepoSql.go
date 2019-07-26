// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

import (
	"bitbucket.org/reidev/restlib/users"
	"database/sql"
	"log"
)

/**
Define a struct for Repo for use with users
*/
type RepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//Also store the table name
	tableName string

	//Store the required statements to reduce comput time
	getUserRoles   *sql.Stmt
	clearUserRoles *sql.Stmt
	addUserRole    *sql.Stmt

	//We need the role Repo
	permTable PermissionTable
}

//Provide a method to make a new UserRepoSql
func NewRepoMySql(db *sql.DB, tableName string, roleRepo PermissionTable) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		tableName: tableName,
		permTable: roleRepo,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id int NOT NULL AUTO_INCREMENT, userId int, roleId int, PRIMARY KEY (id) )")
	if err != nil {
		log.Fatal(err)
	}

	//Add calc data to table
	getRoles, err := db.Prepare("SELECT roleId FROM " + tableName + " WHERE userId = ? ")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.getUserRoles = getRoles

	//Clear all roles of a user
	clearRoles, err := db.Prepare("DELETE  FROM " + tableName + " WHERE userId = ? ")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.clearUserRoles = clearRoles

	//Clear all roles of a user
	addRole, err := db.Prepare("INSERT INTO " + tableName + "(userId,roleId) VALUES (?, ?)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.addUserRole = addRole

	//Return a point
	return &newRepo

}

//Provide a method to make a new UserRepoSql
func NewRepoPostgresSql(db *sql.DB, tableName string, roleRepo PermissionTable) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		tableName: tableName,
		permTable: roleRepo,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id SERIAL PRIMARY KEY, userId int NOT NULL, roleId int NOT NULL )")
	if err != nil {
		log.Fatal(err)
	}

	//Add calc data to table
	getRoles, err := db.Prepare("SELECT roleId FROM " + tableName + " WHERE userId = $1 ")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.getUserRoles = getRoles

	//Clear all roles of a user
	clearRoles, err := db.Prepare("DELETE  FROM " + tableName + " WHERE userId = $1 ")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.clearUserRoles = clearRoles

	//Clear all roles of a user
	addRole, err := db.Prepare("INSERT INTO " + tableName + "(userId,roleId) VALUES ($1, $2)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.addUserRole = addRole

	//Return a point
	return &newRepo

}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RepoSql) GetPermissions(user users.User) (*Permissions, error) {
	//Get a list of roles
	permissions := make([]string, 0)

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	rows, err := repo.getUserRoles.Query(user.Id())

	//Rows is the result of a query. Its cursor starts before  the first row of the result set. Use Next to advance through the rows:
	defer rows.Close()
	for rows.Next() {
		//Get the role id
		var roleId int
		err = rows.Scan(&roleId)

		//Get the permissions
		rolePermissions := repo.permTable.GetPermissions(roleId)

		//Push back
		permissions = append(permissions, rolePermissions...)

	}
	rows.Close()
	err = rows.Err() // get any error encountered ing iteration

	//If there is an error
	if err != nil {
		return nil, err
	}

	//Get the permissions from
	return &Permissions{
		Permissions: permissions,
	}, nil

}

/**
Get all of the roles
*/
func (repo *RepoSql) GetRoleIds(user users.User) ([]int, error) {
	//Get a list of roles
	roles := make([]int, 0)

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	rows, err := repo.getUserRoles.Query(user.Id())

	//Build the list

	//Rows is the result of a query. Its cursor starts before  the first row of the result set. Use Next to advance through the rows:
	defer rows.Close()
	for rows.Next() {
		//Get the role id
		var roleId int
		err = rows.Scan(&roleId)

		//Push back
		roles = append(roles, roleId)

	}
	rows.Close()
	err = rows.Err() // get any error encountered ing iteration

	//If there is an error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RepoSql) SetRolesByRoleId(user users.User, roles []int) error {
	//Get all of the
	currentRoles, err := repo.GetRoleIds(user)

	//If the roles dont' equal replace them
	if err != nil || !sameRoles(currentRoles, roles) {

		//Clear all of the roles
		_, err := repo.clearUserRoles.Exec(user.Id())

		//Now add each role
		for _, roleId := range roles {
			_, err = repo.addUserRole.Exec(user.Id(), roleId)
		}

		//Check for error
		return err
	}
	return nil
}

/**
Set the user's roles.  Note this wipes out all current roles
*/
func (repo *RepoSql) SetRolesByName(user users.User, roles []string) error {
	//Build a list of roles to add
	roleIds := make([]int, 0)

	//March over the roles
	for _, role := range roles {
		//Get the id
		roleId, err := repo.permTable.LookUpRoleId(role)

		//Add it to the list
		if err == nil {
			roleIds = append(roleIds, roleId)
		}
	}

	//Now update the roles
	return repo.SetRolesByRoleId(user, roleIds)
}

/**
Clean up the database, nothing much to do
*/
func (repo *RepoSql) CleanUp() {
	repo.getUserRoles.Close()

}

//See if the roles are the same
func sameRoles(a, b []int) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

//func RepoDestroyCalc(id int) error {
//	for i, t := range usersList {
//		if t.Id == id {
//			usersList = append(usersList[:i], usersList[i+1:]...)
//			return nil
//		}
//	}
//	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
//}
