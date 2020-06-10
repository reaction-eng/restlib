// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/reaction-eng/restlib/configuration"
	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
)

/**
 * This struct is used to get data from the post command
 */
type FacebookLoginToken struct {
	// The user handler needs to have access to user repo
	AccessToken string `json:"accessToken"`
}

/**
 * This struct is used to get data from the post command
 */
type facebookAccessTokenResponse struct {
	// The user handler needs to have access to user repo
	AccessToken string `json:"access_token"`
}

/**
 * This struct is used to get data from the post command
 */
type facebookTokenDebugResponse struct {
	// The user handler needs to have access to user repo
	Data facebookTokenDebugData `json:"data"`
}

/**
 * This struct is used to get data from the post command
 */
type facebookTokenDebugData struct {
	// The user handler needs to have access to user repo
	AppId string `json:"app_id"`
}

/**
* This struct is used to get data from the post command
 */
type facebookMeResponse struct {
	Email string                 `json:"email"`
	Error map[string]interface{} `json:"error"`
}

/**
 * This struct is used
 */
type FacebookHandler struct {
	// The user handler needs to have access to user repo
	helper Helper

	//Store the facebook info
	clientId     string
	clientSecret string
}

/**
 * This struct is used
 */
func NewFacebookHandler(helper Helper, configuration configuration.Configuration) *FacebookHandler {
	//Create a new
	facebook := &FacebookHandler{
		helper:       helper,
		clientId:     configuration.GetString("facebook_client_id"),
		clientSecret: configuration.GetString("facebook_client_secret"),
	}

	return facebook
}

/**
Function used to get routes
*/
func (fbHandler *FacebookHandler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{ //Allow for the user to login
			Name:        "UserLogin Facebook",
			Method:      "POST",
			Pattern:     "/users/login/facebook",
			HandlerFunc: fbHandler.handleUserLoginFacebook,
			Public:      true,
		},
	}

	return routes

}

/**
Get user email from token
*/
func (fbHandler *FacebookHandler) tokenToEmail(token FacebookLoginToken) (string, error) {
	//Now add the client id
	tokenParams := url.Values{}
	tokenParams.Set("client_id", fbHandler.clientId)
	tokenParams.Set("client_secret", fbHandler.clientSecret)
	tokenParams.Set("grant_type", "client_credentials")

	//Now get my access_token
	response, err := http.Get("https://graph.facebook.com/oauth/access_token?" + tokenParams.Encode())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	//Now convert to a FacebookAccessTokenResponse
	myToken := facebookAccessTokenResponse{}

	//Now decode
	if json.NewDecoder(response.Body).Decode(&myToken) != nil {
		return "", err
	}

	//Now make sure the org user token is for this app
	//Now add the client id
	debugParams := url.Values{}
	debugParams.Set("input_token", token.AccessToken)
	debugParams.Set("access_token", myToken.AccessToken)

	//Now get my access_token
	response, err = http.Get("https://graph.facebook.com/debug_token?" + debugParams.Encode())
	if err != nil {
		return "", err
	}

	//Now convert to a FacebookAccessTokenResponse
	tokenCheck := facebookTokenDebugResponse{}

	//Now decode
	if json.NewDecoder(response.Body).Decode(&tokenCheck) != nil {
		return "", err
	}

	//Make sure that the app id equals my id
	if tokenCheck.Data.AppId != fbHandler.clientId {
		return "", errors.New("token_not_valid_for_this_app")
	}

	//Now look up the user
	//Now make sure the org user token is for this app
	//Now add the client id
	meParams := url.Values{}
	meParams.Set("fields", "email")
	meParams.Set("access_token", token.AccessToken)

	//Now get my access_token
	response, err = http.Get("https://graph.facebook.com/me?" + meParams.Encode())
	if err != nil {
		return "", err
	}

	//Now convert to a FacebookAccessTokenResponse
	meToken := facebookMeResponse{}

	//Now decode
	if json.NewDecoder(response.Body).Decode(&meToken) != nil {
		return "", err
	}
	//If there is no email
	//If there is an error message
	if meToken.Error != nil {
		return "", errors.New(fmt.Sprint(meToken.Error["message"]))
	}

	if len(meToken.Email) == 0 {
		return "", errors.New("requires email in facebook permissions")
	}

	return meToken.Email, nil
}

/**
Function used to create new user
*/
func (fbHandler *FacebookHandler) handleUserLoginFacebook(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	cred := FacebookLoginToken{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Get the users email
	email, err := fbHandler.tokenToEmail(cred)

	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now get the user by email
	user, err := fbHandler.helper.GetUserByEmail(email)

	//See if it a new error
	if err != nil && user == nil {
		//The email is not in use, so add it
		//Create an empty new user
		newUser := fbHandler.helper.NewEmptyUser()
		newUser.SetEmail(email)
		newUser.SetPassword("") //This is a blank password that prevents being able to login

		//Now store it
		user, err = fbHandler.helper.AddUser(newUser)

		//Make sure it created an id
		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
		}

		//Now activate user
		fbHandler.helper.ActivateUser(user)

		//Now get the user again
		//Now get the user by email
		user, err = fbHandler.helper.GetUserByEmail(email)

		if err != nil {
			utils.ReturnJsonError(w, http.StatusForbidden, err)
		}

	} else if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//Create JWT token and Store the token in the response
	user.SetToken(fbHandler.helper.CreateJWTToken(user.Id(), -1, user.Email()))

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusCreated, user)
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}

}
