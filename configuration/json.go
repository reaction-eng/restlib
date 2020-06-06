// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Json struct {
	//Load in the params from json
	params map[string]interface{}

	fatal func(interface{})
}

func NewJsonFatal(configFiles ...string) *Json {
	config, err := NewJson(configFiles...)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func NewJson(configFiles ...string) (*Json, error) {
	//Define a Configuration
	config := Json{
		params: make(map[string]interface{}, 0),
		fatal: func(i interface{}) {
			log.Fatal(i)
		},
	}

	//Now march over each file
	for _, configFile := range configFiles {
		//See if itself a config
		var testMap map[string]interface{}
		err := json.Unmarshal([]byte(configFile), &testMap)

		//See if it is a config string
		if err == nil && testMap != nil {
			//Merge the maps
			for k, v := range testMap {
				config.params[k] = v
			}

		} else {
			//Parse as file
			//Load in the file
			configFileStream, err := os.Open(configFile)

			if err == nil {
				//Get the json and add to the params
				jsonParser := json.NewDecoder(configFileStream)
				jsonParser.Decode(&config.params)
				configFileStream.Close()
			}
		}

	}

	//Return it
	return &config, nil
}

func (jsonConfig *Json) Get(key string) interface{} {
	//Get the key from the file
	param := jsonConfig.params[key]

	//Now see if it is specified in the env
	systemEnvVar := os.Getenv(key)

	//If it is not empty set it
	if len(systemEnvVar) != 0 {
		param = systemEnvVar
	}

	return param

}

func (jsonConfig *Json) GetFatal(key string) interface{} {
	//Get the thing
	thing := jsonConfig.Get(key)

	//Make sure it is not nil
	if thing == nil {
		jsonConfig.fatal("Cannot not find configuration for " + key)
	}

	return thing

}

func (jsonConfig *Json) GetString(key string) string {
	//Get the value
	value := jsonConfig.Get(key)

	if value == nil {
		return ""
	} else {
		//Get the key from the
		return fmt.Sprint(value)
	}

}

func (jsonConfig *Json) GetStringError(key string) (string, error) {
	//Get the value
	value := jsonConfig.Get(key)

	if value == nil {
		return "", errors.New("could not find " + key)
	}

	//Get the key from the
	return fmt.Sprint(value), nil

}

func (jsonConfig *Json) GetStringFatal(key string) string {
	//Get the key from the
	return fmt.Sprint(jsonConfig.GetFatal(key))

}

func (jsonConfig *Json) GetInt(key string) (int, error) {
	intString, err := jsonConfig.GetStringError(key)
	if err != nil {
		return 0, err
	}

	res, err := strconv.Atoi(intString)

	return res, err

}

func (jsonConfig *Json) GetIntFatal(key string) int {

	result, err := jsonConfig.GetInt(key)

	if err != nil {
		jsonConfig.fatal("Cannot not find int configuration for " + key)
	}

	return result
}

func (jsonConfig *Json) GetFloat(key string) (float64, error) {
	floatString, err := jsonConfig.GetStringError(key)
	if err != nil {
		return 0, err
	}

	res, err := strconv.ParseFloat(floatString, 64)
	return res, err
}

func (jsonConfig *Json) GetKeys() []string {
	keys := make([]string, len(jsonConfig.params))

	i := 0
	for k := range jsonConfig.params {
		keys[i] = k
		i++
	}

	return keys

}

func (jsonConfig *Json) GetConfig(key string) Configuration {
	//Get the child interface
	childConfigInterface := jsonConfig.Get(key)

	//If childConfigInterface, return nil
	if childConfigInterface == nil {
		return nil
	}

	//Now cast it
	childConfig, isMap := childConfigInterface.(map[string]interface{})
	if !isMap {
		return nil
	}

	return &Json{childConfig, jsonConfig.fatal}

}

func (jsonConfig *Json) GetStruct(key string, object interface{}) error {
	//Get the child interface
	childConfigInterface := jsonConfig.Get(key)

	//If childConfigInterface, return nil
	if childConfigInterface == nil {
		return errors.New("cannot not find int configuration for " + key)
	}

	//Now unmarshal
	jsonByte, err := json.Marshal(childConfigInterface)

	//If there is no error
	if err != nil {
		return err
	}

	//Now put the json back into the object
	err = json.Unmarshal(jsonByte, object)

	return err
}

func (jsonConfig *Json) GetStringArray(key string) []string {
	//Get the child interface
	childConfigInterface := jsonConfig.Get(key)

	//Get as an array
	childArray, isArray := childConfigInterface.([]interface{})
	if !isArray {
		return nil
	}

	//Now build a new slice
	childStringArray := make([]string, 0)

	//Now march over each child array
	for _, child := range childArray {
		childStringArray = append(childStringArray, fmt.Sprint(child))
	}

	return childStringArray

}

// GetBool returns a configuration entry typed as bool.

// key is the configuration entry to retrieve. If the entry does not exist or
// is not bool, then defaultVal is returned. Values of type string that are
// some variant of true/false, True/False, etc. will be converted to bool. This
// also works for parameters retrieved from environment variables.
//
// It is the caller's responsibility to provide the correct default for the
// given use case.
func (jsonConfig *Json) GetBool(key string, defaultVal bool) bool {
	value := jsonConfig.Get(key)
	if value == nil {
		return defaultVal
	}

	if boolValue, isBool := value.(bool); isBool {
		return boolValue
	}

	if stringValue, err := strconv.ParseBool(fmt.Sprint(value)); err == nil {
		return stringValue
	}

	return defaultVal
}
