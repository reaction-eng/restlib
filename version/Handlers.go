package version

import (
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/utils"
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
