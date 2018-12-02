package roles

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

//Simple struct to hold the role
type Role struct {
	Name        string
	Permissions []string
}

/**
Define a struct for Repo for use with users
*/
type PermissionTableJson struct {
	//Hold a map of the strings
	Roles map[int]Role
}

//Provide a method to make a new UserRepoSql
func NewPermissionTableJson(fileName string) *PermissionTableJson {
	//Create a new table
	permTable := &PermissionTableJson{}

	//Load in the file
	configFileStream, err := os.Open(fileName)
	defer configFileStream.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	//Get the json and add to the Params
	jsonParser := json.NewDecoder(configFileStream)
	jsonParser.Decode(&permTable)

	return permTable
}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *PermissionTableJson) GetPermissions(roleId int) []string {
	//Look up the role
	role := repo.Roles[roleId]

	//Else get them
	return role.Permissions

}

/**
Get the role id for this name
*/
func (repo *PermissionTableJson) LookUpRoleId(roleLookUp string) (int, error) {
	//March over each config
	for index, role := range repo.Roles {
		//If the role equals
		if roleLookUp == role.Name {
			return index, nil
		}
	}

	//It was not found, error an error
	return -1, errors.New("could not find role " + roleLookUp)

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
