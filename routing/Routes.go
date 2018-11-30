package routing

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name           string
	Method         string
	Pattern        string
	HandlerFunc    http.HandlerFunc
	Public         bool
	ReqPermissions []string
	MuxRoute       *mux.Route
}
