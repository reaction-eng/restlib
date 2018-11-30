package roles

import (
	"bitbucket.org/reidev/restlib/users"
	"database/sql"
	"log"
)

/**
Define a struct for Repo for use with users
*/
type PermissionRepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//Also store the table name
	tableName string

	//Store the required statements to reduce comput time
	getUserRoles *sql.Stmt

	//We need the role Repo
	roleRepo RoleRepo
}

//Provide a method to make a new UserRepoSql
func NewRepoMySql(db *sql.DB, tableName string, roleRepo RoleRepo) *PermissionRepoSql {

	//Define a new repo
	newRepo := PermissionRepoSql{
		db:        db,
		tableName: tableName,
		roleRepo:  roleRepo,
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
	//Store it
	newRepo.getUserRoles = getRoles

	//Return a point
	return &newRepo

}

//Provide a method to make a new UserRepoSql
func NewRepoPostgresSql(db *sql.DB, tableName string, roleRepo RoleRepo) *PermissionRepoSql {

	//Define a new repo
	newRepo := PermissionRepoSql{
		db:        db,
		tableName: tableName,
		roleRepo:  roleRepo,
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
	//Store it
	newRepo.getUserRoles = getRoles

	//Return a point
	return &newRepo

}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *PermissionRepoSql) GetPermissions(user users.User) (*Permissions, error) {
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
		rolePermissions := repo.roleRepo.GetPermissions(roleId)

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
		permissions: permissions,
	}, nil

}

/**
Clean up the database, nothing much to do
*/
func (repo *PermissionRepoSql) CleanUp() {
	repo.getUserRoles.Close()

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
