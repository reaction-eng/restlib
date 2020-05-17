// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
)

type Handler struct {
	//Store the user helper
	userHelper Helper

	//Keep track if we want to allow userCreation
	allowUserCreation bool
}

/**
 * This struct is used
 */
func NewHandler(userHelper Helper, allowUserCreation bool) *Handler {
	//Build a new User Handler
	handler := Handler{
		userHelper:        userHelper,
		allowUserCreation: allowUserCreation,
	}

	return &handler
}

/**
Function used to get routes
*/
func (handler *Handler) GetRoutes() []routing.Route {

	//Provide the user update and documentation by default
	var routes = make([]routing.Route, 0)

	//If the user can create users append the routes
	if handler.allowUserCreation {

		routes = append(routes,
			routing.Route{ //Now for the user info
				Name:        "UserCreate",
				Method:      "POST",
				Pattern:     "/users/new",
				HandlerFunc: handler.handleUserCreate,
				Public:      true,
			},
			routing.Route{ //Allow the user to turn on their account
				Name:        "User Activate",
				Method:      "POST",
				Pattern:     "/users/activate",
				HandlerFunc: handler.handleUserActivationPut,
				Public:      true,
			},
			routing.Route{ //Allow the user to turn on their account
				Name:        "User Activate",
				Method:      "GET",
				Pattern:     "/users/activate",
				HandlerFunc: handler.handleUserActivationGet,
				Public:      true,
			},
			routing.Route{ //Allow the user to turn on their account
				Name:        "Get User Activation Token",
				Method:      "GET",
				Pattern:     "/users/activate",
				HandlerFunc: handler.handleUserActivationPut,
				Public:      true,
			},
			routing.Route{ //Allow for the user to get an update of them selves
				Name:        "PasswordChange",
				Method:      "POST",
				Pattern:     "/users/password/change",
				HandlerFunc: handler.handlePasswordUpdate,
				Public:      false,
			},
			routing.Route{ //Allow for the user to ask for a password change
				Name:        "PasswordResetGet",
				Method:      "GET",
				Pattern:     "/users/password/reset",
				HandlerFunc: handler.handlePasswordResetGet,
				Public:      true,
			},
			routing.Route{ //Allow the user to set their password
				Name:        "PasswordResetPost",
				Method:      "POST",
				Pattern:     "/users/password/reset",
				HandlerFunc: handler.handlePasswordResetPut,
				Public:      true,
			},
		)

	}

	//Add in the normal routes
	routes = append(routes,
		routing.Route{ //Allow for the user to login
			Name:        "UserLogin",
			Method:      "POST",
			Pattern:     "/users/login",
			HandlerFunc: handler.handleUserLogin,
			Public:      true,
		},
		routing.Route{ //Allow for the user to login
			Name:        "User Api Documentation",
			Method:      "GET",
			Pattern:     "/api/users",
			HandlerFunc: handler.handleUserDocumentation,
			Public:      true,
		},
		routing.Route{ //Allow for the user to update them selves
			Name:        "UserUpdate",
			Method:      "PUT",
			Pattern:     "/users/",
			HandlerFunc: handler.handleUserUpdate,
			Public:      false,
		},
		routing.Route{ //Allow for the user to get an update of them selves
			Name:        "UserGet",
			Method:      "GET",
			Pattern:     "/users/",
			HandlerFunc: handler.handleUserGet,
			Public:      false,
		},
	)

	return routes

}

/**
Function used to create new user
*/
func (handler *Handler) handleUserCreate(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	newUser := handler.userHelper.NewEmptyUser()

	type newUserStruct struct {
		Email          string `json:"email"`
		Password       string `json:"password"`
		OrganizationId int    `json:"organizationId`
	}

	//Create the new user
	newUserInfo := &newUserStruct{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(newUserInfo)
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
		return
	}

	//Copy over the new user data
	newUser.SetEmail(newUserInfo.Email)
	newUser.SetPassword(newUserInfo.Password)
	newUser.SetOrganizations(newUserInfo.OrganizationId)

	//Now create the new user
	err = handler.userHelper.CreateUser(newUser)

	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
		return
	}
	utils.ReturnJsonStatus(w, http.StatusCreated, true, "create_user_added")
}

/**
Function used to create new user
*/
func (handler *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type loginUserStruct struct {
		Email          string `json:"email"`
		Password       string `json:"password"`
		OrganizationId int    `json:"organizationId"`
	}

	userCred := &loginUserStruct{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(userCred)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	user, err := handler.userHelper.GetUserByEmail(strings.TrimSpace(strings.ToLower(userCred.Email)))

	//check for an error
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//We have the user, try to login
	user, err = handler.userHelper.Login(userCred.Password, userCred.OrganizationId, user)

	//If there is an error, don't login
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	utils.ReturnJson(w, http.StatusCreated, user)
}

/**
updates everything but the password for the user
*/
func (handler *Handler) handleUserUpdate(w http.ResponseWriter, r *http.Request) {

	loggedInUserString := r.Context().Value("user") //Grab the id of the user that send the request
	if loggedInUserString == nil {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("no_user_logged_in"))
		return
	}

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := loggedInUserString.(int) //Grab the id of the user that send the request

	//Now load the current user from the repo
	user, err := handler.userHelper.GetUser(loggedInUser)

	//Check for an error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//decode the request body into struct with all of the info specified and failed if any error occur
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now update the user
	user, err = handler.userHelper.Update(loggedInUser, user)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusAccepted, user)
	} else {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
	}

}

/**
Get the current up to date user
*/
func (handler *Handler) handleUserGet(w http.ResponseWriter, r *http.Request) {

	loggedInUserString := r.Context().Value("user") //Grab the id of the user that send the request
	if loggedInUserString == nil {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("no_user_logged_in"))
		return
	}
	loggedInUser := loggedInUserString.(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userHelper.GetUser(loggedInUser)

	//Check to see if the user was created
	if err == nil {
		//Make sure we null the password
		//Blank out the password before returning
		user.SetPassword("")

		utils.ReturnJson(w, http.StatusOK, user)
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnsupportedMediaType, false, err.Error())
	}

}

/**
Updates the password for this user
*/
func (handler *Handler) handlePasswordUpdate(w http.ResponseWriter, r *http.Request) {
	loggedInUserString := r.Context().Value("user") //Grab the id of the user that send the request
	if loggedInUserString == nil {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("no_user_logged_in"))
		return
	}
	loggedInUser := loggedInUserString.(int) //Grab the id of the user that send the request

	//Create a new password change object
	type updatePasswordChangeStruct struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		PasswordOld string `json:"passwordold"`
	}
	info := updatePasswordChangeStruct{}

	//Now get the json info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now update the password
	err = handler.userHelper.PasswordChange(loggedInUser, info.Email, info.Password, info.PasswordOld)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusAccepted, true, "password_change_success")
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}

}

/**
Function to request a password change
*/
func (handler *Handler) handlePasswordResetGet(w http.ResponseWriter, r *http.Request) {

	//Now get the email that was passed in
	keys, ok := r.URL.Query()["email"]

	//Only take the first one
	if !ok || len(keys[0]) < 1 {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, "password_change_missing_email")
		return
	}

	//Get the email
	email := keys[0]

	//Look up the user
	user, err := handler.userHelper.GetUserByEmail(email)

	//If there is an error just return, we don't want people to know if there was an email here
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusOK, true, "password_change_request_received")
		return
	}

	//Now issue a request
	err = handler.userHelper.IssueResetRequest(handler.userHelper.TokenGenerator(), user.Id(), user.Email())

	//There was a real error return
	if err != nil {
		utils.ReturnJsonError(w, http.StatusServiceUnavailable, err)
		return
	}

	//Now just return
	utils.ReturnJsonStatus(w, http.StatusOK, true, "password_change_request_received")

}

func (handler *Handler) handlePasswordResetPut(w http.ResponseWriter, r *http.Request) {

	//Define a local struct to get the email out of the request
	type ResetGet struct {
		Email      string `json:"email"`
		ResetToken string `json:"reset_token"`
		Password   string `json:"password"`
	}

	//Create a new password change object
	info := ResetGet{}

	//Now get the json info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Lookup the user id
	user, err := handler.userHelper.GetUserByEmail(info.Email)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "password_change_forbidden")
		return
	}

	//Try to use the token
	requestId, err := handler.userHelper.CheckForResetToken(user.Id(), info.ResetToken)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "password_change_forbidden")
		return
	}

	//Now update the password
	err = handler.userHelper.PasswordChangeForced(user.Id(), user.Email(), info.Password)
	//Return the error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}
	//Mark the request as used
	err = handler.userHelper.UseToken(requestId)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusAccepted, true, "password_change_success")
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}
}

func (handler *Handler) handleUserActivationPut(w http.ResponseWriter, r *http.Request) {

	//Define a local struct to get the email out of the request
	type ActivationGet struct {
		Email    string `json:"email"`
		ActToken string `json:"activation_token"`
	}

	//Create a new password change object
	info := ActivationGet{}

	//Now get the json info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Lookup the user id
	user, err := handler.userHelper.GetUserByEmail(info.Email)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "activation_forbidden")
		return
	}

	//Try to use the token
	requestId, err := handler.userHelper.CheckForActivationToken(user.Id(), info.ActToken)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "activation_forbidden")
		return
	}
	//Now activate the user
	err = handler.userHelper.ActivateUser(user)

	//Return the error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}
	//Mark the request as used
	err = handler.userHelper.UseToken(requestId)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusAccepted, true, "user_activated")
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}
}

/**
Get a new user activation token
*/
func (handler *Handler) handleUserActivationGet(w http.ResponseWriter, r *http.Request) {

	//Now get the email that was passed in
	keys, ok := r.URL.Query()["email"]

	//Only take the first one
	if !ok || len(keys[0]) < 1 {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, "activation_token_missing_email")
		return
	}

	//Get the email
	email := keys[0]

	//Look up the user
	user, err := handler.userHelper.GetUserByEmail(email)

	//If there is an error just return, we don't want people to know if there was an email here
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusOK, true, "activation_token_request_received")
		return
	}

	//Now issue a request
	//If the user is not already active
	if user.Activated() {
		utils.ReturnJsonStatus(w, http.StatusOK, true, "activation_token_request_received")
		return
	}
	//Else issue the request
	err = handler.userHelper.IssueActivationRequest(handler.userHelper.TokenGenerator(), user.Id(), user.Email())

	if err != nil {
		utils.ReturnJsonError(w, http.StatusServiceUnavailable, err)
		return
	}

	//Now just return
	utils.ReturnJsonStatus(w, http.StatusOK, true, "activation_token_request_received")
}
