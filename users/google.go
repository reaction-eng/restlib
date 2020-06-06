// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/reaction-eng/restlib/configuration"
	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2"
	goauth2 "google.golang.org/api/oauth2/v2"
)

/**
 * This struct is used to get data from the post command
 */
type GoogleLoginToken struct {
	// The user handler needs to have access to user repo
	IdToken string `json:"id_token"`
}

/**
 * This struct is used
 */
type GoogleHandler struct {
	// The user handler needs to have access to user repo
	helper Helper

	//We need the oauth config
	oAuthConfig *oauth2.Config
}

/**
 * This struct is used
 */
func NewGoogleHandler(helper Helper, configuration configuration.Configuration) *GoogleHandler {
	//Create a new
	google := &GoogleHandler{
		helper: helper,
		oAuthConfig: &oauth2.Config{
			ClientID:     configuration.GetStringFatal("google_client_id"),
			ClientSecret: configuration.GetStringFatal("google_client_secret"),
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}

	return google
}

/**
Function used to get routes
*/
func (gHandler *GoogleHandler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{ //Allow for the user to login
			Name:        "UserLogin Google",
			Method:      "POST",
			Pattern:     "/users/login/google",
			HandlerFunc: gHandler.handleUserLoginGoogle,
			Public:      true,
		},
	}

	return routes

}

/**
Function used to create new user
*/
func (gHandler *GoogleHandler) handleUserLoginGoogle(w http.ResponseWriter, r *http.Request) {

	googleLoginStruct := struct {
		Token          oauth2.Token `json:"token"`
		OrganizationId int          `json:"organizationId"`
	}{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(&googleLoginStruct)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}
	//Make sure it is valid
	if !googleLoginStruct.Token.Valid() {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, errors.New("invalid_token"))
		return

	}

	//Now get the user info
	ctx := context.Background()
	client := oauth2.NewClient(ctx, gHandler.oAuthConfig.TokenSource(ctx, &googleLoginStruct.Token))
	svc, err := goauth2.New(client)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//And get the user info
	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Make sure there is an email
	if len(userInfo.Email) == 0 {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, errors.New("invalid_email"))
		return
	}

	//Now get the user by email
	user, err := gHandler.helper.GetUserByEmail(userInfo.Email)

	//See if it a new error
	if err != nil && user == nil {
		//The email is not in use, so add it
		//Create an empty new user
		newUser := gHandler.helper.NewEmptyUser()
		newUser.SetEmail(userInfo.Email)
		newUser.SetPassword("") //This is a blank password that prevents being able to login

		//Now store it
		user, err = gHandler.helper.AddUser(newUser)

		//Make sure it created an id
		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
		}

		//Now activate user
		gHandler.helper.ActivateUser(user)

		//Now get the user again
		//Now get the user by email
		user, err = gHandler.helper.GetUserByEmail(user.Email())

		// add the user to the
		gHandler.helper.AddUserToOrganization(user, googleLoginStruct.OrganizationId)

		user, err = gHandler.helper.GetUserByEmail(user.Email())

		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
		}

	} else if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	if !InOrganization(user, googleLoginStruct.OrganizationId) {
		utils.ReturnJsonError(w, http.StatusForbidden, errors.New("not in organization"))
		return
	}

	//Create JWT token and Store the token in the response
	user.SetToken(gHandler.helper.CreateJWTToken(user.Id(), googleLoginStruct.OrganizationId, user.Email()))

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusCreated, user)
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}
}
