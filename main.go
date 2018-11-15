package restlib

import (
	"bitbucket.org/reidev/restlib/authentication"
	"bitbucket.org/reidev/restlib/configuration"
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"bitbucket.org/reidev/stoichcalcservicego/calcs"
	"database/sql"
	"log"
	"net/http"
)

//Define the global variables that are setup in the main
var calcsRepo calcs.Repo

func main() {

	//Get the config so we can build the database
	config := configuration.NewConfiguration()

	//Define using a database
	//Load in the database
	db, err := sql.Open("mysql", config.GetDataBaseSourceName())
	defer db.Close()

	//Check for an error
	if err != nil {
		log.Fatal(err)
	}

	//Define the repos
	calcsRepo = calcs.NewRepoSql(db, "calcs")

	//Defer cleaning up
	defer calcsRepo.CleanUp()

	// We will have users so create a user repo
	userRepo := users.NewRepoSql(db, "users")
	defer userRepo.CleanUp()

	//We also need to handle requests about users,
	userHandler := users.NewHandler(userRepo)

	//Define the router, by in the routes specific to this project, and others
	router := routing.NewRouter(nil, []routing.Route{}, userHandler)

	//Add in middleware/filter that respons to CORS
	router.Use(authentication.MakeCORSMiddlewareFunc()) //Make sure to add the cross site permission first

	//Add in middleware/filter that checks for user authentication
	router.Use(authentication.MakeJwtMiddlewareFunc(router, userRepo))

	//Start the filter
	log.Fatal(http.ListenAndServe(config.GetString("host_port"), router))

	//http.HandleFunc("/", DefaultHandle)
	//log.Fatal(http.ListenAndServe(config.GetString("host_port"), nil))

}
