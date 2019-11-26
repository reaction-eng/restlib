// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package utils

import (
	"encoding/json"
	"net/http"
)

/**
Provide a support method to return json
*/
func ReturnJsonMessage(w http.ResponseWriter, statusCode int, message string) {

	//Now just pass it
	ReturnJson(w, statusCode, map[string]interface{}{"message": message})

}

/**
Provide a support method to return json
*/
func ReturnJsonStatus(w http.ResponseWriter, statusCode int, status bool, message string) {

	//Now just pass it
	ReturnJson(w, statusCode, map[string]interface{}{"status": status, "message": message})

}

/**
Provide a support method to return json
*/
func ReturnJsonError(w http.ResponseWriter, statusCode int, err error) {

	//Now just pass it
	ReturnJson(w, statusCode, map[string]interface{}{"status": false, "message": err.Error()})

}

/**
Provide a support method to return json
*/
func ReturnJson(w http.ResponseWriter, statusCode int, data interface{}) {

	//Assume it is always json
	w.Header().Set("Content-GetType", "application/json; charset=UTF-8")

	//Pass in the code
	w.WriteHeader(statusCode) // unprocessable entity

	//Now encode the json object
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}

/**
Provide a support method to return json
*/
func ReturnJsonNoEscape(w http.ResponseWriter, statusCode int, data interface{}) {

	//Assume it is always json
	w.Header().Set("Content-GetType", "application/json; charset=UTF-8")

	//Pass in the code
	w.WriteHeader(statusCode) // unprocessable entity

	//Setup a new encoder
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	//Now encode the json object
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}

}
