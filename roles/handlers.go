// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package roles

import (
	"net/http"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/users"
	"github.com/reaction-eng/restlib/utils"
)

type Handler struct {
	// The user handler needs to have access to user repo
	userRepo users.Repo

	//Store the repo for the roles
	roleRepo Repo
}

func NewHandler(userRepo users.Repo, roleRepo Repo) *Handler {
	//Build a new User Handler
	handler := Handler{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}

	return &handler
}

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

func (handler *Handler) handleUserPermissionsGet(w http.ResponseWriter, r *http.Request) {

	//We have gone through the auth, so we should know the id of the logged in user
	loggedInUser, organizationId, err := utils.UserFromContext(r.Context())
	if err != nil {
		utils.ReturnJsonError(w, http.StatusForbidden, err)
		return
	}

	//Get the user
	user, err := handler.userRepo.GetUser(loggedInUser)

	//If there is no error
	if err != nil {
		utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Get the list of permissions
	perm, err := handler.roleRepo.GetPermissions(user, organizationId)

	//Check to see if the user was created
	if err == nil {
		utils.ReturnJson(w, http.StatusOK, perm)
	} else {
		utils.ReturnJsonStatus(w, http.StatusServiceUnavailable, false, err.Error())
	}
}
