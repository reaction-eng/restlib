// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package static

import (
	"github.com/gorilla/mux"
	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
	"html/template"
	"net/http"
)

/**
Return the routes for this
*/
func GetRoutes(staticRepo Repo) []routing.Route {
	return []routing.Route{
		{ ////////////////////////////////////////////////////////////////////////////////////////////////////////
			Name:    "Help Api Documentation",
			Method:  "GET",
			Pattern: "/api/static",
			Public:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				//Load int he welcome html
				tmpl := template.Must(template.ParseFiles("documentation/api.static.html"))

				//Show it
				tmpl.Execute(w, nil)
			},
		},
		{ ////////////////////////////////////////////////////////////////////////////////////////////////////////
			Name:    "Get Public Page",
			Method:  "GET",
			Pattern: "/static/public/{path}",
			Public:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				path := vars["path"]

				//Get that specific calc
				html, err := staticRepo.GetStaticPublicDocument(path)

				//Check for an error
				if err != nil {
					utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
					return
				}

				//Now show it
				utils.ReturnJson(w, http.StatusOK, html)
			},
		},
		{ ////////////////////////////////////////////////////////////////////////////////////////////////////////
			Name:    "Get Private Page",
			Method:  "GET",
			Pattern: "/static/private/{path}",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				path := vars["path"]

				//Get that specific calc
				html, err := staticRepo.GetStaticPrivateDocument(path)

				//Check for an error
				if err != nil {
					utils.ReturnJsonError(w, http.StatusUnprocessableEntity, err)
					return
				}

				//Now show it
				utils.ReturnJson(w, http.StatusOK, html)
			},
		},
	}

}
