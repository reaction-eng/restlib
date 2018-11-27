package authentication

import (
	"database/sql"
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
	addResetStatement *sql.Stmt
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

	//Add reset data to table
	addReset, err := db.Prepare("INSERT INTO " + tableName + "(userId,email, token, issued) VALUES (?, ?, ?, ?)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	//Store it
	newRepo.addResetStatement = addReset

	//Return a point
	return &newRepo

}

/**
Look up the user and return if they were found
*/
func (repo *PasswordResetRepoSql) issueResetRequest(userId int, email string) error {

	//Get a new token
	token := TokenGenerator()

	//Now add it to the database
	//Add the info
	//execute the statement//(userId,name,input,flow)- "(userId,email, token, issued)
	_, err := repo.addResetStatement.Exec(userId, email, token, time.Now().Format(time.RFC3339))

	//Now email
	//TODO: Email
	fmt.Println(userId, email, token)

	//Return the user calcs
	return err
}

/**
Clean up the database, nothing much to do
*/
func (repo *PasswordResetRepoSql) CleanUp() {
	repo.addResetStatement.Close()

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
