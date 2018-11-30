package roles

import (
	"bitbucket.org/reidev/restlib/configuration"
	"strconv"
)

/**
Define a struct for Repo for use with users
*/
type RoleRepoJson struct {
	//Hold on to the sql databased
	db *configuration.Configuration
}

//Provide a method to make a new UserRepoSql
func NewRoleRepoJson(fileName string) *RoleRepoJson {
	//Load in the json
	db := configuration.NewConfiguration(fileName)

	//Add that to the role repo
	return &RoleRepoJson{
		db: db,
	}

}

/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RoleRepoJson) GetPermissions(roleId int) []string {
	//Look up the role
	role := repo.db.GetConfig(strconv.Itoa(roleId))

	//If not nil
	if role == nil {
		return []string{}
	}

	//Else get them
	return role.GetStringArray("permissions")
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
