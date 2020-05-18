package mocks

import (
	"context"
	"net/http"
	"strings"

	"github.com/reaction-eng/restlib/users"

	"github.com/reaction-eng/restlib/utils"

	"github.com/reaction-eng/restlib/routing"
)

type TestRouter struct {
	routes []routing.Route
	userId *int
	orgId  *int
}

func NewTestRouter(routerProducer routing.RouteProducer) *TestRouter {
	return &TestRouter{
		routerProducer.GetRoutes(),
		nil,
		nil,
	}
}
func NewTestRouterWithUserId(routerProducer routing.RouteProducer, userId int, orgId int) *TestRouter {
	return &TestRouter{
		routerProducer.GetRoutes(),
		&userId,
		&orgId,
	}
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
