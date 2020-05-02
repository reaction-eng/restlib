// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"database/sql"

	"github.com/reaction-eng/restlib/users"
)

const TableName = "userpref"

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

func NewRepoMySql(db *sql.DB, baseOptions *OptionGroup) (*RepoSql, error) {
	//Define a new repo
	newRepo := RepoSql{
		db:          db,
		baseOptions: baseOptions,
	}

	getSetting, err := db.Prepare("SELECT settings FROM " + TableName + " WHERE userID = ?")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.getSettingFromDbCmd = getSetting

	setSetting, err := db.Prepare("INSERT INTO " + TableName + "(userId,settings) VALUES (?,?) ON DUPLICATE KEY UPDATE settings = VALUES(settings)")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.setSettingIntoDbCmd = setSetting

	return &newRepo, nil
}

func NewRepoPostgresSql(db *sql.DB, baseOptions *OptionGroup) (*RepoSql, error) {
	//Define a new repo
	newRepo := RepoSql{
		db:          db,
		baseOptions: baseOptions,
	}

	getSetting, err := db.Prepare("SELECT settings FROM " + TableName + " WHERE userID = $1")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.getSettingFromDbCmd = getSetting

	//Get the settings
	setSetting, err := db.Prepare("INSERT INTO " + TableName + "(userId,settings) VALUES ($1, $2) ON CONFLICT (userId) DO UPDATE SET settings = $2")
	//Check for error
	if err != nil {
		return nil, err
	}
	newRepo.setSettingIntoDbCmd = setSetting

	return &newRepo, nil
}

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

func (repo *RepoSql) CleanUp() {
	//Close all of the prepared statements
	repo.getSettingFromDbCmd.Close()
	repo.setSettingIntoDbCmd.Close()
}
