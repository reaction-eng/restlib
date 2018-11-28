package authentication

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

/**
Define a struct for Repo for use with users
*/
type PasswordResetRepoSql struct {
	//Hold on to the sql databased
	db *sql.DB

	//Also store the table name
	tableName string

	//We need the emailer
	emailer PasswordResetEmailer

	//Store the required statements to reduce comput time
	addRequestStatement *sql.Stmt
	getRequestStatement *sql.Stmt
	rmRequestStatement  *sql.Stmt
}

//Provide a method to make a new UserRepoSql
func NewRepoSql(db *sql.DB, tableName string, emailer PasswordResetEmailer) *PasswordResetRepoSql {

	//Define a new repo
	newRepo := PasswordResetRepoSql{
		db:        db,
		tableName: tableName,
		emailer:   emailer,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id int NOT NULL AUTO_INCREMENT, userId int, email TEXT, token TEXT, issued DATE, PRIMARY KEY (id) )")
	if err != nil {
		log.Fatal(err)
	}

	//Add request data to table
	addRequest, err := db.Prepare("INSERT INTO " + tableName + "(userId,email, token, issued) VALUES (?, ?, ?, ?)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.addRequestStatement = addRequest

	//pull the request from the table
	getRequest, err := db.Prepare("SELECT * FROM " + tableName + " where userId = ? AND token = ?")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.getRequestStatement = getRequest

	//pull the request from the table
	rmRequest, err := db.Prepare("delete FROM " + tableName + " where id = ? limit 1")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.rmRequestStatement = rmRequest

	//Return a point
	return &newRepo

}

/**
Look up the user and return if they were found
*/
func (repo *PasswordResetRepoSql) IssueResetRequest(userId int, email string) error {

	//Get a new token
	token := TokenGenerator()

	//Now add it to the database
	//Add the info
	//execute the statement//(userId,name,input,flow)- "(userId,email, token, issued)
	_, err := repo.addRequestStatement.Exec(userId, email, token, time.Now())

	//Now email
	//TODO: Email
	fmt.Println(userId, " ", email, " ", token)

	//Return the user calcs
	return err
}

/**
Use the taken to validate
*/
func (repo *PasswordResetRepoSql) CheckForResetToken(userId int, token string) (int, error) {

	//Prepare to get values
	//id,  userId int, email TEXT, token TEXT, issued DATE,
	var id int
	var userIdDB int
	var emailDB string
	var tokenDB string
	var issued time.Time

	//Get the value
	err := repo.getRequestStatement.QueryRow(userId, token).Scan(&id, &userIdDB, &emailDB, &tokenDB, &issued)

	//If there is an error, assume it can't be done
	if err != nil {
		return -1, errors.New("password_change_forbidden")
	}

	//Make sure the user id and token match
	if userId != userIdDB || tokenDB != token {
		return -1, errors.New("password_change_forbidden")
	}

	//So it was correct, check the date
	//TODO: check the date

	//Return the user calcs
	return id, nil
}

func (repo *PasswordResetRepoSql) UseResetToken(id int) error {

	//Remove the token
	_, err := repo.rmRequestStatement.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

/**
Clean up the database, nothing much to do
*/
func (repo *PasswordResetRepoSql) CleanUp() {
	repo.getRequestStatement.Close()
	repo.addRequestStatement.Close()
	repo.rmRequestStatement.Close()

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