// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"bitbucket.org/reidev/restlib/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/**
 * This struct is used
 */
type Handler struct {
	// The user handler needs to have access to user repo
	userRepo users.Repo

	//Store the repo for the roles
	roleRepo Repo
}

/**
 * This struct is used
 */
func NewHandler(userRepo users.Repo, roleRepo Repo) *Handler {
	//Build a new User Handler
	handler := Handler{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}

	return &handler
}

/**
Function used to get routes
*/
func (handler *Handler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{ //Allow for the user to login
			Name:        "Get the User Preferences",
			Method:      "GET",
			Pattern:     "/users/preferences",
			HandlerFunc: handler.handleUserPreferencesGet,
		},
		{ //Allow for the user to login
			Name:        "Set the User Preferences",
			Method:      "POST",
			Pattern:     "/users/preferences",
			HandlerFunc: handler.handleUserPreferencesSet,
		},
	}

	return routes

}

/**
*Get the current up to date user
 */
func (handler *Handler) handleUserPreferencesGet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//If there is no error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Get the list of permissions
	perf, err := handler.roleRepo.GetPreferences(user)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, perf)
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnsupportedMediaType, false, err.Error())
	}

}

/**
*Get the current up to date user
 */
func (handler *Handler) handleUserPreferencesSet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//If there is no error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Load in a limited amount of data from the body
	body, err := ioutil.ReadAll(r.Body)

	//Create an empty new calc
	settings := SettingGroup{}

	//Now marshal it into the body
	if err := json.Unmarshal(body, &settings); err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, err.Error())
		return
	}

	//Get the list of permissions
	pref, err := handler.roleRepo.SetPreferences(user, &settings)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, pref)
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnsupportedMediaType, false, err.Error())
	}

}
