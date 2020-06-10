// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package preferences

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/users"
	"github.com/reaction-eng/restlib/utils"
)

type Handler struct {
	// The user handler needs to have access to user repo
	userRepo users.Repo

	prefRepo Repo
}

func NewHandler(userRepo users.Repo, prefRepo Repo) *Handler {
	//Build a new User Handler
	handler := Handler{
		userRepo: userRepo,
		prefRepo: prefRepo,
	}

	return &handler
}

func (handler *Handler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{
			Name:        "Get the User Preferences",
			Method:      "GET",
			Pattern:     "/users/preferences",
			HandlerFunc: handler.handleUserPreferencesGet,
		},
		{
			Name:        "Set the User Preferences",
			Method:      "POST",
			Pattern:     "/users/preferences",
			HandlerFunc: handler.handleUserPreferencesSet,
		},
	}

	return routes

}

func (handler *Handler) handleUserPreferencesGet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUserString := r.Context().Value("user")
	if loggedInUserString == nil {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("no_user_logged_in"))
		return
	}

	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//If there is no error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	perf, err := handler.prefRepo.GetPreferences(user)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, perf)
	} else {
		utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, err.Error())
	}

}

func (handler *Handler) handleUserPreferencesSet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUserString := r.Context().Value("user") //Grab the id of the user that send the request
	if loggedInUserString == nil {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("no_user_logged_in"))
		return
	}
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//If there is no error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	//Create an empty new calc
	settings := SettingGroup{}

	//Now marshal it into the body
	if err := json.Unmarshal(body, &settings); err != nil {
		utils.ReturnJsonStatus(w, http.StatusBadRequest, false, err.Error())
		return
	}

	pref, err := handler.prefRepo.SetPreferences(user, &settings)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, pref)
	} else {
		utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, err.Error())
	}

}
