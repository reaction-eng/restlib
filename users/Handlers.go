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
	}

	return routes

}

/**
Function used to create new user
*/
func (handler *Handler) handleUserCreate(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	newUser := &User{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
		return

	}

	//Now create the new suer
	err = newUser.Create(handler.userRepo) //Create account

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJsonStatus(w, http.StatusCreated, true, "user "+newUser.Email+" added ")
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnprocessableEntity, false, err.Error())
	}

}

/**
Function used to create new user
*/
func (handler *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request) {

	//Create an empty new user
	userCred := &User{}

	//decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(userCred)
	if err != nil {
		utils.ReturnJson(w, http.StatusUnprocessableEntity, err)
		return

	}

	//Now look up the user
	user, err := handler.userRepo.GetUserByEmail(userCred.Email)

	//check for an error
	if err != nil {
		//There prob is not a user to return
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//We have the user, try to login
	err = user.Login(userCred.Password)

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
