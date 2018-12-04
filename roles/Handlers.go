package roles

import (
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"bitbucket.org/reidev/restlib/utils"
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
			Name:        "User Api Documentation",
			Method:      "GET",
			Pattern:     "/api/users/permissions",
			HandlerFunc: handler.handlePermissionsDocumentation,
			Public:      true,
		},
		{ //Allow for the user to login
			Name:        "Get the User Permissions",
			Method:      "GET",
			Pattern:     "/users/permissions",
			HandlerFunc: handler.handleUserPermissionsGet,
		},
	}

	return routes

}

/**
*Get the current up to date user
 */
func (handler *Handler) handleUserPermissionsGet(w http.ResponseWriter, r *http.Request) {

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
	perm, err := handler.roleRepo.GetPermissions(user)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, perm)
	} else {
		utils.ReturnJsonStatus(w, http.StatusUnsupportedMediaType, false, err.Error())
	}

}
