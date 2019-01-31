package users

import (
	"bitbucket.org/reidev/restlib/utils"
	"database/sql"
	"errors"
	"log"
	"time"
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
	addUserStatement        *sql.Stmt
	getUserStatement        *sql.Stmt
	getUserByEmailStatement *sql.Stmt
	updateUserStatement     *sql.Stmt
	activateStatement       *sql.Stmt
	listAllUsersStatement   *sql.Stmt

	//Store the nullable time object

}

//Provide a method to make a new UserRepoSql
func NewRepoMySql(db *sql.DB, tableName string) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		tableName: tableName,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, activation Date, PRIMARY KEY (id) )")
	if err != nil {
		log.Fatal(err)
	}

	//Add calc data to table
	addUser, err := db.Prepare("INSERT INTO " + tableName + "(email,password) VALUES (?, ?)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.addUserStatement = addUser

	//get user statement
	getUser, err := db.Prepare("SELECT * FROM " + tableName + " where id = ?")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.getUserStatement = getUser

	//get calc statement
	getUserByEmail, err := db.Prepare("SELECT * FROM " + tableName + " where email like ?")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.getUserByEmailStatement = getUserByEmail

	//update the user
	updateStatement, err := db.Prepare("UPDATE  " + tableName + " SET email = ?, password = ? WHERE id = ?")

	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.updateUserStatement = updateStatement

	//Activate User statemetn
	activateStatement, err := db.Prepare("UPDATE  " + tableName + " SET activation = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.activateStatement = activateStatement

	//update the user
	listAllUsers, err := db.Prepare("SELECT id FROM " + tableName)
	if err != nil {
		log.Fatal(err)
	}
	newRepo.listAllUsersStatement = listAllUsers

	//Return a point
	return &newRepo

}

//Provide a method to make a new UserRepoSql
func NewRepoPostgresSql(db *sql.DB, tableName string) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		tableName: tableName,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id SERIAL PRIMARY KEY, email TEXT NOT NULL, password TEXT NOT NULL, activation Date)")
	if err != nil {
		log.Fatal(err)
	}

	//Add calc data to table
	addUser, err := db.Prepare("INSERT INTO " + tableName + "(email,password) VALUES ($1, $2)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	////Store it
	newRepo.addUserStatement = addUser

	//get calc statement
	getUser, err := db.Prepare("SELECT * FROM " + tableName + " where id = $1")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	////Store it
	newRepo.getUserStatement = getUser

	//get calc statement
	getUserByEmail, err := db.Prepare("SELECT * FROM " + tableName + " where email like $1")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.getUserByEmailStatement = getUserByEmail

	//update the user
	updateStatement, err := db.Prepare("UPDATE  " + tableName + " SET email = $1, password = $2 WHERE id = $3")

	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.updateUserStatement = updateStatement

	//Activate User statemetn
	activateStatement, err := db.Prepare("UPDATE  " + tableName + " SET activation = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.activateStatement = activateStatement

	//update the user
	listAllUsers, err := db.Prepare("SELECT id FROM " + tableName)
	if err != nil {
		log.Fatal(err)
	}
	newRepo.listAllUsersStatement = listAllUsers

	//Return a point
	return &newRepo

}

/**
Look up the user and return if they were found
*/
func (repo *RepoSql) GetUserByEmail(email string) (User, error) {
	//var dataResult string
	var user BasicUser

	//Store the sql time
	var activationDate utils.NullTime

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserByEmailStatement.QueryRow(email).Scan(&user.Id_, &user.Email_, &user.password_, &activationDate)

	//Use a useful error
	if err == sql.ErrNoRows {
		err = errors.New("login_email_not_found")
		return nil, err
	}

	//Store if this is activated
	user.activated_ = activationDate.Valid
	user.passwordlogin_ = len(user.password_) > 0

	//Return the user calcs
	return &user, err
}

/**
Look up the user by id and return if they were found
*/
func (repo *RepoSql) GetUser(id int) (User, error) {
	//var dataResult string
	var user BasicUser

	//Store the sql time
	var activationDate utils.NullTime

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserStatement.QueryRow(id).Scan(&user.Id_, &user.Email_, &user.password_, &activationDate)

	//Use a useful error
	if err == sql.ErrNoRows {
		err = errors.New("login_user_id_not_found")
	}

	//Store if this is activated
	user.activated_ = activationDate.Valid
	user.passwordlogin_ = len(user.password_) > 0

	//Return the user calcs
	return &user, err
}

/**
List all of the users
*/
func (repo *RepoSql) ListAllUsers() ([]int, error) {
	//Put in the list
	list := make([]int, 0)

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	rows, err := repo.listAllUsersStatement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		//Append the row
		list = append(list, id)

	}
	err = rows.Err()

	return list, err
}

/**
Add the user to the database
*/
func (repo *RepoSql) AddUser(newUser User) (User, error) {

	//Add the info
	//execute the statement//(userId,name,input,flow)
	_, err := repo.addUserStatement.Exec(newUser.Email(), newUser.Password())

	//Check for error
	if err != nil {
		return newUser, err
	}

	//Now look up the person by email
	return repo.GetUserByEmail(newUser.Email())

}

/**
Update the user table.  No checks are made here,
*/
func (repo *RepoSql) UpdateUser(user User) (User, error) {
	//Update the user statement
	//Just update the info
	//execute the statement//"UPDATE  " + tableName + " SET email = ?, password = ? WHERE id = ?"
	_, err := repo.updateUserStatement.Exec(user.Email(), user.Password(), user.Id())

	//Check for error
	if err != nil {
		log.Fatal(err)
	}

	return user, err
}

/**
Update the user table.  No checks are made here,
*/
func (repo *RepoSql) ActivateUser(user User) error {
	//Get the current time
	actTime := utils.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	//Just update the info//"UPDATE  " + tableName + " SET activation = $1 WHERE id = $2")
	_, err := repo.activateStatement.Exec(actTime, user.Id())

	//Check for error
	if err != nil {
		log.Fatal(err)
	}

	return err
}

/**
Clean up the database, nothing much to do
*/
func (repo *RepoSql) CleanUp() {
	repo.addUserStatement.Close()
	repo.getUserByEmailStatement.Close()
	repo.getUserStatement.Close()
	repo.updateUserStatement.Close()
	repo.listAllUsersStatement.Close()
}

/**
Clean up the database, nothing much to do
*/
func (repo *RepoSql) NewEmptyUser() User {
	return &BasicUser{}
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
