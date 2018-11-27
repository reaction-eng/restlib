package users

import (
	"database/sql"
	"errors"
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
	addUserStatement        *sql.Stmt
	getUserStatement        *sql.Stmt
	getUserByEmailStatement *sql.Stmt
}

//Provide a method to make a new UserRepoSql
func NewRepoSql(db *sql.DB, tableName string) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:        db,
		tableName: tableName,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id) )")
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

	//get calc statement
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

	//Return a point
	return &newRepo

}

/**
Look up the user and return if they were found
*/
func (repo *RepoSql) GetUserByEmail(email string) (User, error) {
	//var dataResult string
	var user BasicUser

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserByEmailStatement.QueryRow(email).Scan(&user.Id_, &user.Email_, &user.Password_)

	//Use a useful error
	if err == sql.ErrNoRows {
		err = errors.New("login_email_not_found")
	}

	//Return the user calcs
	return &user, err
}

/**
Look up the user by id and return if they were found
*/
func (repo *RepoSql) GetUser(id int) (User, error) {
	//var dataResult string
	var user BasicUser

	//Get the value //id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, PRIMARY KEY (id)
	err := repo.getUserStatement.QueryRow(id).Scan(&user.Id_, &user.Email_, &user.Password_)

	//Use a useful error
	if err == sql.ErrNoRows {
		err = errors.New("login_user_id_not_found")
	}

	//Return the user calcs
	return &user, err
}

/**
Add the user to the database
*/
func (repo *RepoSql) AddUser(newUser User) (User, error) {
	//Add the info
	//execute the statement//(userId,name,input,flow)
	result, err := repo.addUserStatement.Exec(newUser.Email(), newUser.Password())

	//Check for error
	if err != nil {
		return newUser, err
	}

	////Get the id
	newId, _ := result.LastInsertId()

	//Add the newid to the user and return it
	newUser.SetId(int(newId))

	return newUser, nil
}

/**
Clean up the database, nothing much to do
*/
func (repo *RepoSql) CleanUp() {
	repo.addUserStatement.Close()
	repo.getUserByEmailStatement.Close()
	repo.getUserStatement.Close()
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
