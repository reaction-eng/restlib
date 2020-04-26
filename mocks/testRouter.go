package mocks

import (
	"context"
	"net/http"
	"strings"

	"github.com/reaction-eng/restlib/routing"
	"github.com/reaction-eng/restlib/users"
)

type TestRouter struct {
	routes []routing.Route
	user   users.User
}

func NewTestRouter(routerProducer routing.RouteProducer, user users.User) *TestRouter {
	return &TestRouter{
		routerProducer.GetRoutes(),
		user,
	}
}

func (router *TestRouter) Handle(w http.ResponseWriter, r *http.Request) *routing.Route {
	// Update the context
	if router.user != nil {
		ctx := context.WithValue(r.Context(), "user", router.user.Id())
		r = r.WithContext(ctx)
	}

	uri := r.URL.Path
	method := r.Method

	for _, route := range router.routes {
		if strings.EqualFold(route.Pattern, uri) && strings.EqualFold(route.Method, method) {
			route.HandlerFunc(w, r)

			return &route
		}
	}
	return nil
}
