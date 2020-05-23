// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

import (
	"database/sql"
	"sort"

	"github.com/reaction-eng/restlib/users"
)

const TableName = "roles"

type RepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//Store the required statements to reduce compute time
	getUserRoles   *sql.Stmt
	clearUserRoles *sql.Stmt
	addUserRole    *sql.Stmt

	//We need the role Repo
	permTable PermissionTable
}

func NewRepoMySql(db *sql.DB, roleRepo PermissionTable) (*RepoSql, error) {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		permTable: roleRepo,
	}

	getRoles, err := db.Prepare("SELECT roleId FROM " + TableName + " WHERE userId = ? AND orgId = ?")
	if err != nil {
		return nil, err
	}
	newRepo.getUserRoles = getRoles

	clearRoles, err := db.Prepare("DELETE FROM " + TableName + " WHERE userId = ? AND orgId = ?")
	if err != nil {
		return nil, err
	}
	newRepo.clearUserRoles = clearRoles

	addRole, err := db.Prepare("INSERT INTO " + TableName + " (userId,orgId,roleId) VALUES (?, ?, ?)")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.addUserRole = addRole

	return &newRepo, nil
}

func NewRepoPostgresSql(db *sql.DB, permTable PermissionTable) (*RepoSql, error) {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		permTable: permTable,
	}

	getRoles, err := db.Prepare("SELECT roleId FROM " + TableName + " WHERE userId = $1 AND orgId = $2")
	if err != nil {
		return nil, err
	}
	newRepo.getUserRoles = getRoles

	clearRoles, err := db.Prepare("DELETE  FROM " + TableName + " WHERE userId = $1 AND orgId = $2")
	if err != nil {
		return nil, err
	}
	newRepo.clearUserRoles = clearRoles

	addRole, err := db.Prepare("INSERT INTO " + TableName + " (userId,orgId,roleId) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	newRepo.addUserRole = addRole

	return &newRepo, nil

}

func (repo *RepoSql) GetPermissions(user users.User, organizationId int) (*Permissions, error) {
	//Get a list of roles
	permissions := make([]string, 0)

	rows, err := repo.getUserRoles.Query(user.Id(), organizationId)
	if err != nil {
		return nil, err
	}

	//Rows is the result of a query. Its cursor starts before  the first row of the result set. Use Next to advance through the rows:
	defer rows.Close()
	for rows.Next() {
		//Get the role id
		var roleId int
		err = rows.Scan(&roleId)
		if err != nil {
			return nil, err
		}

		//Get the permissions
		rolePermissions := repo.permTable.GetPermissions(roleId)

		//Push back
		permissions = append(permissions, rolePermissions...)

	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//Get the permissions from
	return &Permissions{
		Permissions: permissions,
	}, nil

}

func (repo *RepoSql) GetRoleIds(user users.User, organizationId int) ([]int, error) {
	//Get a list of roles
	roles := make([]int, 0)

	rows, err := repo.getUserRoles.Query(user.Id(), organizationId)
	if err != nil {
		return nil, err
	}

	//Rows is the result of a query. Its cursor starts before  the first row of the result set. Use Next to advance through the rows:
	defer rows.Close()
	for rows.Next() {
		//Get the role id
		var roleId int
		err = rows.Scan(&roleId)
		if err != nil {
			return nil, err
		}
		//Push back
		roles = append(roles, roleId)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (repo *RepoSql) SetRolesByRoleId(user users.User, organizationId int, roles []int) error {
	//Get all of the
	currentRoles, err := repo.GetRoleIds(user, organizationId)
	userId := user.Id()

	if err != nil {
		return err
	}

	//If the roles dont' equal replace them
	if err != nil || !sameRoles(currentRoles, roles) {

		//Clear all of the roles
		_, err := repo.clearUserRoles.Exec(userId, organizationId)
		if err != nil {
			return err
		}

		//Now add each role
		for _, roleId := range roles {
			_, err = repo.addUserRole.Exec(userId, organizationId, roleId)
			if err != nil {
				return err
			}
		}

		return err
	}
	return nil
}

func (repo *RepoSql) SetRolesByName(user users.User, organizationId int, roles []string) error {
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

	return repo.SetRolesByRoleId(user, organizationId, roleIds)
}

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

	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
