// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type SettingGroup struct {
	Settings map[string]string `json:"settings"`

	SubGroup map[string]*SettingGroup `json:"subgroup"`
}

func newSettingGroup() *SettingGroup {
	return &SettingGroup{
		Settings: make(map[string]string),
		SubGroup: make(map[string]*SettingGroup),
	}
}

/**
Make sure that the settings and subgroup are valid objects
*/
func (setGroup *SettingGroup) checkSubStructureValid() {
	if setGroup.Settings == nil {
		setGroup.Settings = make(map[string]string)
	}
	if setGroup.SubGroup == nil {
		setGroup.SubGroup = make(map[string]*SettingGroup)
	}

}

//Provide a way to get the sub group
func (setGroup *SettingGroup) GetSubGroup(id string) *SettingGroup {
	//See if there is a settings group
	if group, found := setGroup.SubGroup[id]; found {
		return group
	} else {
		//Create the subgroup
		newSubGroup := newSettingGroup()

		//Store it setGroup
		setGroup.SubGroup[id] = newSubGroup

		return newSubGroup
	}
}

func (setGroup *SettingGroup) GetValueAsString(id string) (string, error) {
	//See if there is a settings group
	if group, found := setGroup.Settings[id]; found {
		return group, nil
	} else {
		return "", errors.New("setting " + id + " not found")
	}
}

func (setGroup *SettingGroup) GetValueAsBool(id string) (bool, error) {
	//Get the value
	value, err := setGroup.GetValueAsString(id)
	if err != nil {
		return false, err
	}

	//Now convert to bool
	valueBool, _ := strconv.ParseBool(value)
	return valueBool, nil

}

func (setGroup *SettingGroup) GetSettingAsBool(tree []string) (bool, error) {
	//Start with the current subGroup
	subGroup := setGroup

	//March down the list of subgroups until we get the tree
	for i := 0; i < len(tree)-1; i++ {
		subGroup = subGroup.GetSubGroup(tree[i])

	}

	//Now get the last value
	return subGroup.GetValueAsBool(tree[len(tree)-1])
}

//Provide a way to get the sub group
func (setGroup *SettingGroup) GetSettingAsString(tree []string) (string, error) {
	//Start with the current subGroup
	subGroup := setGroup

	//March down the list of subgroups until we get the tree
	for i := 0; i < len(tree)-1; i++ {
		subGroup = subGroup.GetSubGroup(tree[i])
	}

	//Now get the last value
	return subGroup.GetValueAsString(tree[len(tree)-1])
}

//Provide a way to get the sub group
func (setGroup *SettingGroup) checkAndSetDefaultValues(options *OptionGroup) {

	//Now march over and set each value
	for _, opt := range options.Options {
		//See if there is a settings group
		if _, found := setGroup.Settings[opt.Id]; !found {
			setGroup.Settings[opt.Id] = opt.DefaultValue
		}
	}

	//Now do this for all of the subgroups
	for _, optGroup := range options.SubGroups {
		//Get the setSubGroup
		setSubGroup := setGroup.GetSubGroup(optGroup.Id)

		//Now update the subgruop
		setSubGroup.checkAndSetDefaultValues(&optGroup)
	}

}

/**
Define custom methods to serialize and un serialize for sql
*/
func (setGroup SettingGroup) Value() (driver.Value, error) {
	//Convert to a string as json
	jsonByte, err := json.Marshal(setGroup)

	//If there is an error return it
	if err != nil {
		return nil, err
	}

	//convert to string
	jsonString := string(jsonByte)

	//Now return
	return driver.Value(jsonString), nil
}

// Implements sql.Scanner. Simplistic -- only handles string and []byte
func (setGroup *SettingGroup) Scan(src interface{}) error {

	//Get the byte
	var source []byte

	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("incompatible type for SettingGroup")
	}

	//If there is some sources
	if len(source) > 0 {
		//Now unmarshal
		return json.Unmarshal(source, &setGroup)

	}
	return nil

}
