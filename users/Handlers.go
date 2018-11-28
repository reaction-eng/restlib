package users

import (
	"bitbucket.org/reidev/restlib/authentication"
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/utils"
	"encoding/json"
	"net/http"
)

/**
 * This struct is used
 */
type Handler struct {
	// The user handler needs to have access to user repo
	userRepo Repo

	//passwordResetRepo
	resetRepo authentication.PasswordResetRepo
}

/**
 * This struct is used
 */
func NewHandler(userRepo Repo, resetRepo authentication.PasswordResetRepo) *Handler {
	//Build a new User Handler
	handler := Handler{
		userRepo:  userRepo,
		resetRepo: resetRepo,
	}

	return &handler
}

/**
Function used to get routes
*/
func (handler *Handler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
		{ //Allow for the user to login
			Name:        "User Api Documentation",
			Method:      "GET",
			Pattern:     "/api/users",
			HandlerFunc: handler.handleUserDocumentation,
			Public:      true,
		},
		{ //Now for the user info
			Name:        "UserCreate",
			Method:      "POST",
			Pattern:     "/users/new",
			HandlerFunc: handler.handleUserCreate,
			Public:      true,
		},
		{ //Allow for the user to login
			Name:        "UserLogin",
			Method:      "POST",
			Pattern:     "/users/login",
			HandlerFunc: handler.handleUserLogin,
			Public:      true,
		},
		{ //Allow for the user to update tthem selves
			Name:        "UserUpdate",
			Method:      "PUT",
			Pattern:     "/users/",
			HandlerFunc: handler.handleUserUpdate,
			Public:      false,
		},
		{ //Allow for the user to get an update of them selves
			Name:        "UserGet",
			Method:      "GET",
			Pattern:     "/users/",
			HandlerFunc: handler.handleUserGet,
			Public:      false,
		},
		{ //Allow for the user to get an update of them selves
			Name:        "PasswordChange",
			Method:      "POST",
			Pattern:     "/users/password/change",
			HandlerFunc: handler.handlePasswordUpdate,
			Public:      false,
		},
		{ //Allow for the user to ask for a password change
			Name:        "PasswordResetGet",
			Method:      "GET",
			Pattern:     "/users/password/reset",
			HandlerFunc: handler.handlePasswordResetGet,
			Public:      true,
		},
		{ //Allow the user to set their password
			Name:        "PasswordResetPost",
			Method:      "POST",
			Pattern:     "/users/password/reset",
			HandlerFunc: handler.handlePasswordResetPut,
			Public:      true,
		},
	}

	return routes

}

/**
Function used to create new user
*/
func (handler *Handler) handleUserCreate(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	newUser := handler.userRepo.NewEmptyUser()

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
		return

	}

	//Now create the new suer
	_, err = createUser(handler.userRepo, newUser)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusCreated, true, "create_user_added")
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
	}

}

/**
Function used to create new user
*/
func (handler *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	userCred := handler.userRepo.NewEmptyUser()

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(userCred)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now look up the user
	user, err := handler.userRepo.GetUserByEmail(userCred.Email())

	//check for an error
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//We have the user, try to login
	user, err = login(userCred.Password(), user)

	//If there is an error, don't login
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusCreated, user)
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}

}

/**
Updates the password for this user
*/
func (handler *Handler) handleUserUpdate(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Now load the current user from the repo
	user, err := handler.userRepo.GetUser(loggedInUser)

	//Check for an error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//decode the request body into struct with all of the info specified and failed if any error occur
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now update the user
	user, err = updateUser(handler.userRepo, loggedInUser, user)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusAccepted, user)
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}

}

/**
Get the current up to date user
*/
func (handler *Handler) handleUserGet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, user)
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnsupportedMediaType, false, err.Error())
	}

}

/**
Updates the password for this user
*/
func (handler *Handler) handlePasswordUpdate(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser := r.Context().Value("user").(int) //Grab the id of the user that send the request

	//Create a new password change object
	info := updatePasswordChangeStruct{}

	//Now get the json info
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now update the password
	err = passwordChange(handler.userRepo, loggedInUser, info)

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
	}

	//Get the email
	email := keys[0]

	//Look up the user
	user, err := handler.userRepo.GetUserByEmail(email)

	//If there is an error just return, we don't want people to know if there was an email here
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusOK, true, "password_change_request_received")
		return
	}

	//Now issue a request
	err = handler.resetRepo.IssueResetRequest(user.Id(), user.Email())

	//There was a real error return
	if err != nil {
		utils.ReturnJsonError(w, http.StatusNotFound, err)
		return
	}

	//Now just return
	utils.ReturnJsonStatus(w, http.StatusOK, true, "password_change_request_received")

}

/**
Function to request a password change
*/
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
	user, err := handler.userRepo.GetUserByEmail(info.Email)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "password_change_forbidden")
		return
	}

	//Try to use the token
	requestId, err := handler.resetRepo.CheckForResetToken(user.Id(), info.ResetToken)

	//Return the error
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusForbidden, false, "password_change_forbidden")
		return
	}

	//Now update the password
	err = passwordChangeForced(handler.userRepo, user.Id(), user.Email(), info.Password)
	//Return the error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}
	//Mark the request as used
	err = handler.resetRepo.UseResetToken(requestId)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusAccepted, false, "password_change_success")
	} else {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
	}
}
