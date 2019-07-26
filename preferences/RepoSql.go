// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

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
	getSettingFromDbCmd *sql.Stmt
	setSettingIntoDbCmd *sql.Stmt

	//We need the role Repo
	baseOptions *OptionGroup
}

//Provide a method to make a new UserRepoSql
func NewRepoMySql(db *sql.DB, tableName string, baseOptions *OptionGroup) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:          db,
		tableName:   tableName,
		baseOptions: baseOptions,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(userId int NOT NULL, settings TEXT NOT NULL, PRIMARY KEY (userId) )")
	if err != nil {
		log.Fatal(err)
	}

	//Get the settings
	getSetting, err := db.Prepare("SELECT settings FROM " + tableName + " WHERE userID = ?")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.getSettingFromDbCmd = getSetting

	//Get the settings
	setSetting, err := db.Prepare("INSERT INTO " + tableName + "(userId,settings) VALUES (?,?) ON DUPLICATE KEY UPDATE settings = VALUES(settings)")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.setSettingIntoDbCmd = setSetting

	//Return a point
	return &newRepo

}

//Provide a method to make a new UserRepoSql
func NewRepoPostgresSql(db *sql.DB, tableName string, baseOptions *OptionGroup) *RepoSql {

	//Define a new repo
	newRepo := RepoSql{
		db:          db,
		tableName:   tableName,
		baseOptions: baseOptions,
	}

	//Create the table if it is not already there
	//Create a table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(userId SERIAL PRIMARY KEY, settings TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	//Get the settings
	getSetting, err := db.Prepare("SELECT settings FROM " + tableName + " WHERE userID = $1")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.getSettingFromDbCmd = getSetting

	//Get the settings
	setSetting, err := db.Prepare("INSERT INTO " + tableName + "(userId,settings) VALUES ($1, $2) ON CONFLICT (userId) DO UPDATE SET settings = $2")
	//Check for error
	if err != nil {
		log.Fatal(err)
	}
	newRepo.setSettingIntoDbCmd = setSetting

	//Return a point
	return &newRepo

}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RepoSql) GetPreferences(user users.User) (*Preferences, error) {
	//Get the settings from the db
	settings, err := repo.getSettingsFromDb(user)
	if err != nil {
		return nil, err
	}

	//Make sure that the settings is valid
	settings.checkSubStructureValid()

	//Now go over each option to make sure that the default is there
	settings.checkAndSetDefaultValues(repo.baseOptions)

	return &Preferences{
		Settings: settings,
		Options:  repo.baseOptions,
	}, nil

}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RepoSql) getSettingsFromDb(user users.User) (*SettingGroup, error) {
	//Get the id
	var setting *SettingGroup

	//Pull from the database
	err := repo.getSettingFromDbCmd.QueryRow(user.Id()).Scan(&setting)
	//If there is an error return
	if err == nil {
		return setting, nil
	} else if err == sql.ErrNoRows {
		return newSettingGroup(), nil

	} else {
		return nil, err
	}
}

func (repo *RepoSql) SetPreferences(user users.User, userSetting *SettingGroup) (*Preferences, error) {

	//Now add the //(asmId,type,Date, comments)
	_, err := repo.setSettingIntoDbCmd.Exec(user.Id(), userSetting)

	return &Preferences{
		Settings: userSetting,
		Options:  repo.baseOptions,
	}, err

}

/**
Nothing much to do for the clean up
*/
func (repo *RepoSql) CleanUp() {
	//Close all of the prepared statements
	repo.getSettingFromDbCmd.Close()

}
