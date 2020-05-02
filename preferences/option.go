// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"encoding/json"
	"os"
)

type OptionType string

const (
	Int    OptionType = "int"
	String OptionType = "string"
	Float  OptionType = "float"
	Bool   OptionType = "bool"
)

//Restore an options group from a file
func LoadOptionsGroup(jsonFile string) (*OptionGroup, error) {

	optGroup := &OptionGroup{}

	//Load in the file
	configFileStream, err := os.Open(jsonFile)

	if err == nil {
		//Get the json and add to the params
		jsonParser := json.NewDecoder(configFileStream)
		err = jsonParser.Decode(&optGroup)
		if err != nil {
			return nil, err
		}

		configFileStream.Close()
	} else {
		return nil, err
	}

	return optGroup, nil
}

type Option struct {
	//Store the name of th option
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//Store the default value
	DefaultValue string `json:"defaultValue"`

	//Store the type
	Type OptionType `json:"type"`

	//Store min max types if possible
	MaxValue float64 `json:"maxValue,omitempty"`
	MinValue float64 `json:"minValue,omitempty"`

	//Store a list of options
	Selection []string `json:"selection,omitempty"`

	//Set if it is a hidden setting
	Hidden bool `json:"hidden"`
}

/**
Simply stores a group of options for easy display
*/
type OptionGroup struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//Store a list of options
	Options []Option `json:"options"`

	//We can also old other groups
	SubGroups []OptionGroup `json:"subgroups"`
}
