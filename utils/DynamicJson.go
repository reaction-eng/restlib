package utils

import (
	"fmt"
)

type DynamicJson map[string]interface{}

func NewDynamicJson() DynamicJson {

	value := make(map[string]interface{}, 0)

	return value

}

func WrapDynamicJson(input interface{}) DynamicJson {

	//See if it is an map
	asMap, isMap := input.(DynamicJson)

	//If it is a map
	if isMap {
		return asMap
	} else {
		//See if it is map
		if asMapInt, isMapInt := input.(map[string]interface{}); isMapInt {
			return asMapInt
		} else {
			return nil

		}

	}
	return nil

}

/**
Returns the existing or new child object
*/
func (dyJs DynamicJson) GetObject(name string) DynamicJson {

	//Check to see if it is already there
	if obj, ok := dyJs[name]; ok {

		//See if it is an map
		asMap, isMap := obj.(DynamicJson)

		//If it is a map
		if isMap {
			return asMap
		} else {
			//See if it is map
			if asMapInt, isMapInt := obj.(map[string]interface{}); isMapInt {
				return asMapInt
			} else {
				return nil

			}

		}

	} else {
		//Create it
		newObj := NewDynamicJson()

		//Store it
		dyJs[name] = newObj

		//Return the new map
		return newObj

	}
}

/**
Return array
*/
func (dyJs DynamicJson) AppendArray(name string, value interface{}) DynamicJson {
	//Check to see if it is already there
	if obj, ok := dyJs[name]; ok {
		//See if it is an map
		asArray, isMap := obj.([]interface{})

		//If it is a map
		if isMap {
			asArray = append(asArray, value)

			//Store it back
			dyJs[name] = asArray
		}

	} else {
		//Create it
		asArray := make([]interface{}, 0)

		asArray = append(asArray, value)

		//Store it back
		dyJs[name] = asArray

	}

	return dyJs
}

/**
Returns the existing or new child object
*/
func (dyJs DynamicJson) GetValue(name string) string {
	//See if it is three
	if obj, ok := dyJs[name]; ok {
		return fmt.Sprint(obj)
	} else {
		return ""
	}
}

/**
Returns the existing or new child object
*/
func (dyJs DynamicJson) GetNumber(name string) float64 {
	//See if it is three
	if obj, ok := dyJs[name]; ok {

		if v, isInt := obj.(float64); isInt {
			return v
		} else {
			return 0
		}
	} else {
		return 0
	}

}

/**
Returns the existing or new child object
*/
func (dyJs DynamicJson) SetValue(name string, value interface{}) DynamicJson {
	//See if it is three
	dyJs[name] = value

	return dyJs
}

func (dyJs DynamicJson) GetArray(name string) *[]interface{} {
	//Check to see if it is already there
	if obj, ok := dyJs[name]; ok {
		//See if it is an map
		asArray, isArray := obj.([]interface{})

		//If it is a map
		if isArray {
			return &asArray
		} else {
			return nil
		}

	} else {
		//Create it
		asArray := make([]interface{}, 0)

		//Store it back
		dyJs[name] = asArray

		return &asArray

	}

}
