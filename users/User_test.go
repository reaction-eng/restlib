package users_test

import (
	"bitbucket.org/reidev/restlib/middleware"
	"bitbucket.org/reidev/restlib/passwords"
	"bitbucket.org/reidev/restlib/routing"
	"bitbucket.org/reidev/restlib/users"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
Function to carray along the
*/
type routingEnv struct {
	router *routing.Router
}

/**
Perform the testing
*/
func TestUserRoutes(t *testing.T) {

	//Define the list of routes we testing
	var routes = []struct {
		method       string
		path         string
		expectedCode int
	}{ //Now define with
		{"GET", "/api/users", http.StatusOK},
		{"PUT", "/users/", http.StatusForbidden},
	}

	//Now run over each test as a logged out user
	for _, rr := range routes {
		//Get the default env
		env := getDefaultEnv(t)

		//Now run the test
		t.Run("logged out "+rr.path, func(t *testing.T) {

			//In the test function build the request
			req, err := http.NewRequest(rr.method, rr.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rec := httptest.NewRecorder()

			//Get the router from the env and do the action
			env.router.ServeHTTP(rec, req)

			//Make sure the status is correct
			if rec.Result().StatusCode != rr.expectedCode {
				t.Errorf("recived status code %d, expected %d", rec.Result().StatusCode, rr.expectedCode)

			}
		})
	}
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.

	//// Check the status code is what we expect.
	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		status, http.StatusOK)
	//}
	//
	//// Check the response body is what we expect.
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		rr.Body.String(), expected)
	//}

}

/**
Builds the default routing env
*/
func getDefaultEnv(t *testing.T) *routingEnv {

	//Build a config string
	configString := "{\"token_password\": \"RvUP*b7fj9JPJ0*OQ9FlCW%Gg7vNTJWfvV7aQf@u9gWuYQ!S@e9SegAYjh!G%V7btMuGC8g29$qOw\"}"

	//Define a memory repo
	userRepo := users.NewRepoMemory()

	//Make a basic
	passHelper := passwords.NewBasicHelper(configString)

	//Make a user helper
	helper := users.NewUserHelper(userRepo, nil, passHelper)

	//Add some default users
	userOne := users.BasicUser{}
	userOne.SetEmail("one@example.com")
	userOne.SetPassword(passHelper.HashPassword("123456"))
	_, err := userRepo.AddUser(&userOne)

	userTwo := users.BasicUser{}
	userTwo.SetEmail("two@example.com")
	userTwo.SetPassword(passHelper.HashPassword("789012"))
	_, err = userRepo.AddUser(&userTwo)

	if err != nil {
		t.Error(err)
	}

	//Define a new router repo
	//We also need to handle requests about users,
	userHandler := users.NewHandler(helper, false)

	//Define the router, by in the routes specific to this project, and others
	router := routing.NewRouter(nil, nil, nil, userHandler)

	//Add in middleware/filter that respons to CORS
	router.Use(middleware.MakeCORSMiddlewareFunc()) //Make sure to add the cross site permission first

	//Add in middleware/filter that checks for user passwords
	router.Use(middleware.MakeJwtMiddlewareFunc(router, userRepo, nil, passHelper))

	//Define the routing env
	env := routingEnv{
		router: router,
	}

	return &env

}
