// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package version

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/utils"
)

func GetRoutes() []routing.Route {
	return []routing.Route{
		{
			Name:    "Server Version Information",
			Method:  "GET",
			Pattern: "/version",
			Public:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				// Read a simple version file that must be independently generated
				ver, err := ioutil.ReadFile("version.json")
				if err != nil {
					log.Print("Couldn't read version file: ", err)
					utils.ReturnJsonError(w, http.StatusInternalServerError, err)
				}

				w.Write(ver)
			},
		},
	}
}
