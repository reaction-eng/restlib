package mocks

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/reaction-eng/restlib/users"

	"github.com/reaction-eng/restlib/utils"

	"github.com/reaction-eng/restlib/routing"
)

type TestRouter struct {
	router   *mux.Router
	routeMap map[string]routing.Route
	userId   *int
	orgId    *int
}

func NewTestRouter(routerProducer routing.RouteProducer) *TestRouter {
	return NewTestRouterWithUserId(routerProducer, -1, -1)
}
func NewTestRouterWithUserId(routerProducer routing.RouteProducer, userId int, orgId int) *TestRouter {

	var testRouter *TestRouter
	if userId >= 0 {
		testRouter = &TestRouter{
			mux.NewRouter().StrictSlash(true),
			make(map[string]routing.Route, 0),
			&userId,
			&orgId,
		}
	} else {
		// not logged in
		testRouter = &TestRouter{
			mux.NewRouter().StrictSlash(true),
			make(map[string]routing.Route, 0),
			nil,
			nil,
		}
	}

	for _, route := range routerProducer.GetRoutes() {
		muxRoute := testRouter.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

		testRouter.routeMap[muxRoute.GetName()] = route
	}

	return testRouter

}

func NewTestRouterWithUser(routerProducer routing.RouteProducer, user users.User, orgId int) *TestRouter {
	if user == nil {
		return NewTestRouter(routerProducer)
	} else {
		return NewTestRouterWithUserId(routerProducer, user.Id(), orgId)
	}
}

func (router *TestRouter) Handle(w http.ResponseWriter, r *http.Request) *routing.Route {
	// Update the context
	if router.userId != nil {
		ctx := context.WithValue(r.Context(), utils.UserKey, *router.userId)
		ctx = context.WithValue(ctx, utils.OrganizationKey, *router.orgId)
		r = r.WithContext(ctx)
	}

	router.router.ServeHTTP(w, r)

	var match mux.RouteMatch
	if router.router.Match(r, &match) {
		route := router.routeMap[match.Route.GetName()]
		return &route
	}

	return nil
}
