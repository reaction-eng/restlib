package routing

import (
	"net/http"
)

type Route struct {
	Name           string
	Method         string
	Pattern        string
	HandlerFunc    http.HandlerFunc
	Public         bool
	ReqPermissions []string
}
