// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Role struct {
	Name        string
	Permissions []string
}

type PermissionTableJson struct {
	Roles map[int]Role
}

func NewPermissionTableJson(fileName string) (*PermissionTableJson, error) {
	//Create a new table
	permTable := &PermissionTableJson{}

	//Load in the file
	configFileStream, err := os.Open(fileName)
	defer configFileStream.Close()
	if err != nil {
		return nil, err
	}
	//Get the json and add to the params
	jsonParser := json.NewDecoder(configFileStream)
	err = jsonParser.Decode(&permTable)

	if err != nil {
		return nil, err
	}

	return permTable, nil
}

func (repo *PermissionTableJson) GetPermissions(roleId int) []string {
	//Look up the role
	role, hasRole := repo.Roles[roleId]

	if hasRole {
		return role.Permissions
	} else {
		return []string{}
	}
}

func (repo *PermissionTableJson) LookUpRoleId(roleLookUp string) (int, error) {
	//March over each config
	for index, role := range repo.Roles {
		//If the role equals
		if strings.EqualFold(roleLookUp, role.Name) {
			return index, nil
		}
	}

	//It was not found, error an error
	return -1, errors.New("could not find role " + roleLookUp)
}
