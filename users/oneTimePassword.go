// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
)

type OneTimePasswordHandler struct {
	// The user handler needs to have access to user repo
	helper Helper
}

func NewOneTimePasswordHandler(helper Helper) *OneTimePasswordHandler {
	//Create a new
	oneTimePasswordHandler := &OneTimePasswordHandler{
		helper: helper,
	}

	return oneTimePasswordHandler
}

func (handler *OneTimePasswordHandler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{
			Name:        "Get OneTimePassword Token",
			Method:      "GET",
			Pattern:     "/users/onetimelogin",
			HandlerFunc: handler.handleOneTimePasswordGet,
			Public:      true,
		},
		{
			Name:        "Login With OneTimePassword Token",
			Method:      "POST",
			Pattern:     "/users/onetimelogin",
			HandlerFunc: handler.handleOneTimePasswordLoginPut,
			Public:      true,
		},
	}

	return routes

}

func (handler *OneTimePasswordHandler) handleOneTimePasswordGet(w http.ResponseWriter, r *http.Request) {

	//Now get the email that was passed in
	emailKeys, ok := r.URL.Query()["email"]

	//Only take the first one
	if !ok || len(emailKeys[0]) < 1 {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, "onetimepassword_token_missing_email")
		return
	}

	orgKeys, ok := r.URL.Query()["organizationId"]
	if !ok || len(emailKeys[0]) < 1 {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, "onetimepassword_token_missing_organizationId")
		return
	}

	//Get the email
	email := emailKeys[0]
	organizationIdString := orgKeys[0]
	organizationId, err := strconv.Atoi(organizationIdString)
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, "onetimepassword_token_missing_organizationId")
		return
	}

	//Look up the user
	user, err := handler.helper.GetUserByEmail(email)

	//If there is no user create them
	if err == UserNotFound {
		//The email is not in use, so add it
		//Create an empty new user
		newUser := handler.helper.NewEmptyUser()
		newUser.SetEmail(email)
		newUser.SetPassword("") //This is a blank password that prevents being able to login
		newUser.SetOrganizations(organizationId)

		//Now store it
		user, err = handler.helper.AddUser(newUser)
		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
			return
		}
		//Add the users to the org
		for _, orgId := range user.Organizations() {
			err = handler.helper.AddUserToOrganization(newUser, orgId)
			if err != nil {
				utils.ReturnJsonError(w, http.StatusForbidden, err)
				return
			}
		}

		//Make sure it created an id
		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
			return
		}
	} else if err != nil {
		utils.ReturnJsonError(w, http.StatusServiceUnavailable, err)
		return
	}

	err = handler.helper.IssueOneTimePasswordRequest(handler.helper.TokenGenerator(), user.Id(), user.Email())

	if err != nil {
		utils.ReturnJsonError(w, http.StatusServiceUnavailable, err)
		return
	}

	//Now just return
	utils.ReturnJsonStatus(w, http.StatusOK, true, "onetimepassword_token_request_received")
}

func (handler *OneTimePasswordHandler) handleOneTimePasswordLoginPut(w http.ResponseWriter, r *http.Request) {

	//Define a local struct to get the email out of the request
	type LoginGet struct {
		Email          string `json:"email"`
		LoginToken     string `json:"login_token"`
		OrganizationId int    `json:"organizationId`
	}

	info := LoginGet{}

	//Now get the json info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Lookup the user id
	user, err := handler.helper.GetUserByEmail(info.Email)

	//Return the error
	if err == UserNotFound {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "onetimepassword_forbidden")
		return
	} else if err != nil {
		utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, "onetimepassword_forbidden")
		return
	}

	if !InOrganization(user, info.OrganizationId) {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, UserNotInOrganization.Error())
		return
	}

	//Try to use the token
	requestId, err := handler.helper.CheckForOneTimePasswordToken(user.Id(), info.LoginToken)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "onetimepassword_forbidden")
		return
	}

	err = handler.helper.UseToken(requestId)
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//If the user was not activated, activate them
	if !user.Activated() {
		err = handler.helper.ActivateUser(user)
		if err != nil {
			utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, "onetimepassword_forbidden")
			return
		}

		user, err = handler.helper.GetUser(user.Id())
		if err != nil {
			utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, "onetimepassword_forbidden")
			return
		}
	}

	//Create JWT token and Store the token in the response
	user.SetToken(handler.helper.CreateJWTToken(user.Id(), info.OrganizationId, user.Email()))

	utils.ReturnJson(w, http.StatusCreated, user)
}
