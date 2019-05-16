package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
Define a simple database configuration
*/
type Configuration struct {
	//Load in the Params from json
	Params map[string]interface{}
}

//Provide a function to create a new one
func NewConfiguration(configFiles ...string) (*Configuration, error) {
	//Define a Configuration
	config := Configuration{
		Params: make(map[string]interface{}, 0),
	}

	// Read secrets last which will overwrite any existing keys
	configFiles = append(configFiles, "config.secret.json")

	//Now march over each file
	for _, configFile := range configFiles {
		//See if itself a config
		var testMap map[string]interface{}
		err := json.Unmarshal([]byte(configFile), &testMap)

		//See if it is a config string
		if err == nil && testMap != nil {
			//Merge the maps
			for k, v := range testMap {
				config.Params[k] = v
			}

		} else {
			//Parse as file
			//Load in the file
			configFileStream, err := os.Open(configFile)

			if err == nil {
				//Get the json and add to the Params
				jsonParser := json.NewDecoder(configFileStream)
				jsonParser.Decode(&config.Params)
				configFileStream.Close()
			}
		}

	}

	//Return it
	return &config, nil
}

/**
 * Add function to get item
 */
func (config *Configuration) Get(key string) interface{} {
	//Get the key from the file
	param := config.Params[key]

	//Now see if it is specified in the env
	systemEnvVar := os.Getenv(key)

	//If it is not empty set it
	if len(systemEnvVar) != 0 {
		param = systemEnvVar
	}

	return param

}

/**
 * Add function to get item
 */
func (config *Configuration) GetFatal(key string) interface{} {
	//Get the thing
	thing := config.Get(key)

	//Make sure it is not nil
	if thing == nil {
		log.Fatal("Cannot not find configuration for " + key)
	}

	return thing

}

/**
 * Add function to get item
 */
func (config *Configuration) GetString(key string) string {
	//Get the key from the
	return fmt.Sprint(config.Get(key))

}

/**
 * Add function to get item
 */
func (config *Configuration) GetStringError(key string) (string, error) {
	//Get the value
	value := config.Get(key)

	if value == nil {
		return "", errors.New("could not find " + key)
	}

	//Get the key from the
	return fmt.Sprint(value), nil

}

/**
 * Add function to get item
 */
func (config *Configuration) GetStringFatal(key string) string {
	//Get the key from the
	return fmt.Sprint(config.GetFatal(key))

}

/**
 * Add function to get item
 */
func (config *Configuration) GetInt(key string) (int, error) {
	//Get the key from the
	res, err := strconv.Atoi(config.GetString(key))

	return res, err

}

/**
 * Add function to get item
 */
func (config *Configuration) GetIntFatal(key string) int {

	//Get the string
	string := config.GetStringFatal(key)

	//Get the key from the
	res, err := strconv.Atoi(string)

	if err != nil {
		log.Fatal("Cannot not find int configuration for " + key)

	}

	return res

}

/**
 * Add function to get item
 */
func (config *Configuration) GetFloat(key string) (float64, error) {
	//Get the key from the
	res, err := strconv.ParseFloat(config.GetString(key), 64)

	return res, err

}

/**
 * Add function to get item
 */
func (config *Configuration) GetKeys() []string {
	keys := make([]string, len(config.Params))

	i := 0
	for k := range config.Params {
		keys[i] = k
		i++
	}

	return keys

}

/**
 * Add function to get item
 */
func (config *Configuration) GetConfig(key string) *Configuration {
	//Get the child interface
	childConfigInterface := config.Get(key)

	//If childConfigInterface, return nil
	if childConfigInterface == nil {
		return nil
	}

	//Now cast it
	childConfig := childConfigInterface.(map[string]interface{})

	return &Configuration{childConfig}

}

/**
 * Add function to get item
 */
func (config *Configuration) GetStruct(key string, object interface{}) error {
	//Get the child interface
	childConfigInterface := config.Get(key)

	//If childConfigInterface, return nil
	if childConfigInterface == nil {
		return nil
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

/**
 * Add function to get item
 */
func (config *Configuration) GetStringArray(key string) []string {
	//Get the child interface
	childConfigInterface := config.Get(key)

	//Get as an array
	childArray := childConfigInterface.([]interface{})

	//Now build a new slice
	childStringArray := make([]string, 0)

	//Now march over each child array
	for _, child := range childArray {
		childStringArray = append(childStringArray, child.(string))

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
func (config *Configuration) GetBool(key string, defaultVal bool) bool {
	switch val := config.Get(key).(type) {
	case bool:
		return val
	// Convert strings with true/false text to bool
	case string:
		switch strings.ToLower(val) {
		case "true":
			return true
		case "false":
			return false
		}
	}
	return defaultVal
}
