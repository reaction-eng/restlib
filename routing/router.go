package routing

//go:generate mockgen -destination=../mocks/mock_router.go -package=mocks github.com/reaction-eng/restlib/routing Router

import "net/http"

type Router interface {
	GetRoute(req *http.Request) *Route
}
