package users

import (
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
	//passwordResetRepo authentication.
}

/**
 * This struct is used
 */
func NewHandler(userRepo Repo) *Handler {
	//Build a new User Handler
	handler := Handler{
		userRepo: userRepo,
	}

	return &handler
}

/**
Function used to get routes
*/
func (handler *Handler) GetRoutes() []routing.Route {

	var routes = []routing.Route{
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
		{ //Allow for the user to login
			Name:        "User Api Documentation",
			Method:      "GET",
			Pattern:     "/api/users",
			HandlerFunc: handler.handleUserDocumentation,
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
